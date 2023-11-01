package server

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
)

type Doorman struct {
	*pb.UnimplementedDoormanServer

	conn *pgxpool.Pool

	sets    db.Sets
	changes db.Changes
	objects db.Objects
	roles   db.Roles
	tuples  db.Tuples
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
		if i == len(changes)-1 {
			res.PaginationToken = &v.ID
		}
	}
	return res, nil
}

func (d *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	{
		success, err := d.sets.Contains(ctx, doorman.Set{
			Object: doorman.Object(request.Object),
			Verb:   doorman.Verb(request.Verb),
		}, doorman.Object(request.Subject))
		if err != nil && err != db.ErrStale {
			return nil, fmt.Errorf("check failed: %w", err)
		}

		if err == nil {
			return &pb.CheckResponse{Success: success}, nil
		}
	}

	tuples, err := d.tuples.ListParents(ctx, doorman.Object(request.Subject))
	if err != nil {
		return nil, fmt.Errorf("listparents failed: %w", err)
	}

	parentSets, err := doorman.ParentTuplesToSets(ctx, tuples, d.roles.Retrieve)
	if err != nil {
		return nil, fmt.Errorf("pts failed: %w", err)
	}

	_ = d.sets.UpdateParents(ctx, doorman.Object(request.Subject), parentSets)

	// Direct matches first
	for _, tuple := range tuples {
		if tuple.Subject == doorman.Object(request.Subject) && tuple.Object == doorman.Object(request.Object) {
			role, err := d.roles.Retrieve(ctx, tuple.Role)
			if err != nil {
				return nil, fmt.Errorf("db.Retrieve failed: %w", err)
			}
			for _, v := range role.Verbs {
				if v == doorman.Verb(request.Verb) {
					return &pb.CheckResponse{Success: true}, nil
				}
			}
		}
	}

	// Now indirect matches
	for _, tuple := range tuples {
		if tuple.Object.Type() == "group" {
			res, err := d.Check(ctx, &pb.CheckRequest{
				Subject: string(tuple.Object),
				Verb:    request.Verb,
				Object:  request.Object,
			})
			if err != nil {
				return nil, fmt.Errorf("check failed: %w", err)
			}
			if res.Success {
				return &pb.CheckResponse{Success: true}, nil
			}
		}
	}

	return &pb.CheckResponse{Success: false}, nil
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

func (d *Doorman) ListRelations(ctx context.Context, request *pb.ListRelationsRequest) (*pb.ListRelationsResponse, error) {
	return &pb.ListRelationsResponse{}, nil
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

	d.sets.InvalidateParents(ctx, tuple.Subject)

	return &pb.GrantResponse{}, nil
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

	d.sets.InvalidateParents(ctx, tuple.Subject)

	return &pb.RevokeResponse{}, nil
}

func NewDoorman(conn *pgxpool.Pool) *Doorman {
	return &Doorman{conn: conn, changes: db.NewChanges(conn), sets: db.NewSets(conn), roles: db.NewRoles(conn), tuples: db.NewTuples(conn), objects: db.NewObjects(conn)}
}

func mapChangeToPb(c doorman.Change) *pb.Change {
	return &pb.Change{
		Type: c.Type,
	}
}
