package server

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Role:    request.Role,
		Object:  doorman.Object(request.Object),
	}
	if err := d.tuples.Add(ctx, tuple); err != nil {
		return nil, fmt.Errorf("tuples.Add failed: %w", err)
	}

	newTuples := []doorman.Tuple{tuple}
	{
		tupleChildren, err := d.tuples.ListConnected(ctx, tuple.Object, false)
		if err != nil {
			return nil, fmt.Errorf("tuples.ListConnected(obj, false) failed: %w", err)
		}

		tupleParents, err := d.tuples.ListConnected(ctx, tuple.Subject, true)
		if err != nil {
			return nil, fmt.Errorf("tuples.ListConnected(sub, true) failed: %w", err)
		}

		fmt.Println("tup", tuple, tupleChildren, tupleParents)

		// TODO: filter by those that only go through groups

		for _, child := range tupleChildren {
			newTuples = append(newTuples, doorman.Tuple{
				Object:  tuple.Object,
				Role:    child[len(child)-1].Role,
				Subject: child[len(child)-1].Object,
			})
		}

		for _, parent := range tupleParents {
			newTuples = append(newTuples, doorman.Tuple{
				Subject: parent[len(parent)-1].Object,
				Role:    tuple.Role,
				Object:  tuple.Object,
			})
		}

		fmt.Println(newTuples)
		// TODO: connect parents with children??
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
	return nil, status.Errorf(codes.Unimplemented, "method Revoke not implemented")
}

func NewDoorman(relations db.Relations, roles db.Roles, tuples db.Tuples) *Doorman {
	return &Doorman{relations: relations, roles: roles, tuples: tuples}
}
