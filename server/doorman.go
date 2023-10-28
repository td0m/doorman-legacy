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

func (d *Doorman) UpsertRole(ctx context.Context, request *pb.UpsertRoleRequest) (*pb.Role, error) {
	role, err := d.roles.Retrieve(ctx, request.Id)
	if err == db.ErrInvalidRole {
		role = &doorman.Role{ID: request.Id}
	} else if err != nil {
		return nil, fmt.Errorf("db.Retrieve failed: %w", err)
	}

	var tuples []doorman.Tuple
	if role != nil {
		tuples, err = d.tuples.ListTuplesForRole(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("ListTuplesForRole failed: %w", err)
		}
	}

	// THIS SHOULD BE A TX. IF IT FAILS AT ANY POINT AFTER REVOKE THEN WE ARE F*CKED
	// BECAUSE WE HAVE LOST SOURCE OF TRUF

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
		fmt.Println(t)
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

func NewDoorman(relations db.Relations, roles db.Roles, tuples db.Tuples) *Doorman {
	return &Doorman{relations: relations, roles: roles, tuples: tuples}
}
