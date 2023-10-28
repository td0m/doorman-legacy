package server

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	// "golang.org/x/exp/slices"
)

type Doorman struct {
	*pb.UnimplementedDoormanServer

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
	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    doorman.Object(request.Object).Type() + ":" + request.Role,
		Object:  doorman.Object(request.Object),
	}
	newTuples, err := d.tuples.Add(ctx, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuples.Add failed: %w", err)
	}

	relations, err := doorman.TuplesToRelations(ctx, newTuples, d.roles.Retrieve)
	if err != nil {
		return nil, fmt.Errorf("TuplesToRelations failed: %w", err)
	}

	for _, r := range relations {
		if err := d.relations.Add(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to add relation %+v: %w", r, err)
		}
	}

	return &pb.GrantResponse{}, nil
}

func (d *Doorman) Revoke(ctx context.Context, request *pb.RevokeRequest) (*pb.RevokeResponse, error) {
	tuple := doorman.Tuple{
		Subject: doorman.Object(request.Subject),
		Role:    doorman.Object(request.Object).Type() + ":" + request.Role,
		Object:  doorman.Object(request.Object),
	}
	removedTuples, err := d.tuples.Remove(ctx, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuples.Remove failed: %w", err)
	}

	relations, err := doorman.TuplesToRelations(ctx, removedTuples, d.roles.Retrieve)
	if err != nil {
		return nil, fmt.Errorf("TuplesToRelations failed: %w", err)
	}

	for _, r := range relations {
		if err := d.relations.Remove(ctx, r); err != nil {
			return nil, fmt.Errorf("failed to remove relation %+v: %w", r, err)
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

	for _, t := range tuples {
		_, err := d.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w", t, err)
		}
	}

	if err := d.roles.Remove(ctx, request.Id); err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
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

	// TODO: THIS SHOULD BE A TX. IF IT FAILS AT ANY POINT AFTER REVOKE THEN WE ARE F*CKED
	// BECAUSE WE HAVE LOST SOURCE OF TRUF
	// ALT: list connections yourself and compute diff

	for _, t := range tuples {
		_, err := d.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("revoke failed for %s: %w", t, err)
		}
	}

	role.Verbs = []doorman.Verb{}
	for _, v := range request.Verbs {
		role.Verbs = append(role.Verbs, doorman.Verb(v))
	}

	// If this fails all grant requests will fail...
	if err := d.roles.Upsert(ctx, role); err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	for _, t := range tuples {
		_, err := d.Grant(ctx, &pb.GrantRequest{
			Subject: string(t.Subject),
			Role:    doorman.Object(role.ID).Value(),
			Object:  string(t.Object),
		})
		if err != nil {
			return nil, fmt.Errorf("grant failed: %w", err)
		}
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

func NewDoorman(changes db.Changes, relations db.Relations, roles db.Roles, tuples db.Tuples) *Doorman {
	return &Doorman{changes: changes, relations: relations, roles: roles, tuples: tuples}
}
