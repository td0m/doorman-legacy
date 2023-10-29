package server

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	// "golang.org/x/exp/slices"
)

type Doorman struct {
	*pb.UnimplementedDoormanServer

	conn *pgxpool.Pool

	relations db.Relations
	changes   db.Changes
	objects   db.Objects
	roles     db.Roles
	tuples    db.Tuples
}

func (d *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	success, err := d.relations.Check(ctx, doorman.Relation{
		Subject: doorman.Object(request.Subject),
		Verb:    doorman.Verb(request.Verb),
		Object:  doorman.Object(request.Object),
	})
	if err != nil {
		return nil, fmt.Errorf("check failed: %w", err)
	}

	return &pb.CheckResponse{Success: success}, nil
}

func (d *Doorman) Grant(ctx context.Context, request *pb.GrantRequest) (*pb.GrantResponse, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	res, err := d.grantWithTx(ctx, tx, request)
	if err != nil {
		return nil, fmt.Errorf("grantWithTx failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return res, nil
}

func (d *Doorman) grantWithTx(ctx context.Context, tx pgx.Tx, request *pb.GrantRequest) (*pb.GrantResponse, error) {
	if err := d.tuples.WithTx(tx).Lock(ctx); err != nil {
		return nil, fmt.Errorf("tuples.Lock failed: %w, %w", err, tx.Rollback(ctx))
	}

	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    doorman.Object(request.Object).Type() + ":" + request.Role,
		Object:  doorman.Object(request.Object),
	}
	if err := d.tuples.WithTx(tx).Add(ctx, tuple); err != nil {
		return nil, fmt.Errorf("tuples.Add failed: %w, %w", err, tx.Rollback(ctx))
	}

	newTuples, err := d.tuples.WithTx(tx).DependentOn(ctx, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuples.DependentOn failed: %w, %w", err, tx.Rollback(ctx))
	}

	newRelations, err := doorman.TuplesToRelations(ctx, newTuples, d.roles.Retrieve)
	if err != nil {
		return nil, fmt.Errorf("TuplesToRelations failed: %w, %w", err, tx.Rollback(ctx))
	}

	for _, r := range newRelations {
		if err := d.relations.WithTx(tx).Add(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to add relation %+v: %w %w", r, err, tx.Rollback(ctx))
		}
	}

	return &pb.GrantResponse{}, nil
}

func (d *Doorman) Revoke(ctx context.Context, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	res, err := d.revokeWithTx(ctx, tx, request)
	if err != nil {
		return nil, fmt.Errorf("revokeWithTx failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return res, nil
}

func (d *Doorman) revokeWithTx(ctx context.Context, tx pgx.Tx, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	if err := d.tuples.WithTx(tx).Lock(ctx); err != nil {
		return nil, fmt.Errorf("tuples.Lock failed: %w, %w", err, tx.Rollback(ctx))
	}

	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    doorman.Object(request.Object).Type() + ":" + request.Role,
		Object:  doorman.Object(request.Object),
	}
	if err := d.tuples.WithTx(tx).Remove(ctx, tuple); err != nil {
		return nil, fmt.Errorf("tuples.Remove failed: %w, %w", err, tx.Rollback(ctx))
	}

	removedTuples, err := d.tuples.WithTx(tx).DependentOn(ctx, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuples.DependentOn failed: %w, %w", err, tx.Rollback(ctx))
	}

	removedRelations, err := doorman.TuplesToRelations(ctx, removedTuples, d.roles.Retrieve)
	if err != nil {
		return nil, fmt.Errorf("TuplesToRelations failed: %w, %w", err, tx.Rollback(ctx))
	}

	for _, r := range removedRelations {
		if err := d.relations.WithTx(tx).Remove(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to remove relation %+v: %w, %w", r, err, tx.Rollback(ctx))
		}
	}

	return &pb.RevokeResponse{}, nil
}

func (d *Doorman) RemoveRole(ctx context.Context, request *pb.RemoveRoleRequest) (*pb.Role, error) {
	role, err := d.roles.Retrieve(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("db.Retrieve failed: %w", err)
	}

	var tuples []doorman.Tuple
	tuples, err = d.tuples.ListTuplesForRole(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("ListTuplesForRole failed: %w", err)
	}

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	for _, t := range tuples {
		_, err := d.revokeWithTx(ctx, tx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w, %w", t, err, tx.Rollback(ctx))
		}
	}

	if err := d.roles.WithTx(tx).Remove(ctx, request.Id); err != nil {
		return nil, fmt.Errorf("update failed: %w, %w", err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return &pb.Role{}, nil
}

func (d *Doorman) UpsertRole(ctx context.Context, request *pb.UpsertRoleRequest) (*pb.Role, error) {
	role, err := d.roles.Retrieve(ctx, request.Id)
	if err == db.ErrInvalidRole {
		role = &doorman.Role{ID: request.Id}
	} else if err != nil {
		return nil, fmt.Errorf("db.Retrieve failed: %w", err)
	}

	tuples, err := d.tuples.ListTuplesForRole(ctx, role.ID)
	if err != nil {
		return nil, fmt.Errorf("ListTuplesForRole failed: %w", err)
	}

	// THIS SHOULD BE A TX. IF IT FAILS AT ANY POINT AFTER REVOKING AT LEAST 1 THEN WE ARE F*CKED
	// BECAUSE WE HAVE LOST SOURCE OF TRUF
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	for _, t := range tuples {
		_, err := d.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w, %w", t, err, tx.Rollback(ctx))
		}
	}

	role.Verbs = []doorman.Verb{}
	for _, v := range request.Verbs {
		role.Verbs = append(role.Verbs, doorman.Verb(v))
	}

	// If this fails all grant requests will fail...
	if err := d.roles.Upsert(ctx, role); err != nil {
		return nil, fmt.Errorf("update failed: %w, %w", err, tx.Rollback(ctx))
	}

	for _, t := range tuples {
		_, err := d.grantWithTx(ctx, tx, &pb.GrantRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("grant failed: %w, %w", err, tx.Rollback(ctx))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("tx.Commit failed: %w", err)
	}

	return &pb.Role{}, nil
}

func (d *Doorman) ListRelations(ctx context.Context, request *pb.ListRelationsRequest) (*pb.ListRelationsResponse, error) {
	filter := db.RelationFilter{Subject: &request.Subject, Verb: request.Verb}
	relations, err := d.relations.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("db failed: %w", err)
	}

	items := make([]*pb.Relation, len(relations))
	for i, r := range relations {
		items[i] = mapRelationToPb(r)
	}

	return &pb.ListRelationsResponse{
		Items: items,
	}, nil
}

func (d *Doorman) Changes(ctx context.Context, request *pb.ChangesRequest) (*pb.ChangesResponse, error) {
	changes, err := d.changes.List(ctx, db.ChangeFilter{PaginationToken: request.PaginationToken})
	if err != nil {
		return nil, fmt.Errorf("db failed: %w", err)
	}
	res := &pb.ChangesResponse{
		Items: make([]*pb.Change, len(changes)),
	}
	for i, v := range changes {
		res.Items[i] = mapChangeToPb(v)
	}
	return res, nil
}

func mapRelationToPb(r doorman.Relation) *pb.Relation {
	return &pb.Relation{
		Subject: string(r.Subject),
		Verb:    string(r.Verb),
		Object:  string(r.Object),
	}
}

func mapChangeToPb(c doorman.Change) *pb.Change {
	return &pb.Change{
		Type: c.Type,
	}
}

func NewDoorman(conn *pgxpool.Pool) *Doorman {
	return &Doorman{conn: conn, changes: db.NewChanges(conn), relations: db.NewRelations(conn), roles: db.NewRoles(conn), tuples: db.NewTuples(conn), objects: db.NewObjects(conn)}
}
