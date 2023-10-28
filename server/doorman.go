package server

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"golang.org/x/exp/slices"
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
		Role:    doorman.Object(request.Object).Type() + ":" + request.Role,
		Object:  doorman.Object(request.Object),
	}
	if err := d.tuples.Add(ctx, tuple); err != nil {
		return nil, fmt.Errorf("tuples.Add failed: %w", err)
	}

	newTuples := []doorman.TupleWithPath{}
	{
		tupleChildren := []doorman.Path{{}}
		if tuple.Object.Type() == "group" {
			connections, err := d.tuples.ListConnected(ctx, tuple.Object, false)
			if err != nil {
				return nil, fmt.Errorf("tuples.ListConnected(subj, false) failed: %w", err)
			}
			tupleChildren = append(tupleChildren, connections...)
		}

		tupleParents := []doorman.Path{{}}
		if tuple.Subject.Type() == "group" {
			connections, err := d.tuples.ListConnected(ctx, tuple.Subject, true)
			if err != nil {
				return nil, fmt.Errorf("tuples.ListConnected(obj, true) failed: %w", err)
			}
			tupleParents = append(tupleParents, connections...)
		}

		for _, child := range tupleChildren {
			for _, parent := range tupleParents {
				t := doorman.TupleWithPath{
					Tuple: doorman.Tuple{
						Subject: tuple.Subject,
						Role:    tuple.Role,
						Object:  tuple.Object,
					},
					Path: doorman.Path{},
				}
				if len(parent) > 0 {
					t.Subject = parent[len(parent)-1].Object
					path := parent[:len(parent)-1]
					slices.Reverse(path)
					t.Path = append(t.Path, path...)
				}
				if len(child) > 0 {
					t.Object = child[len(child)-1].Object
					t.Role = child[len(child)-1].Role
					t.Path = append(t.Path, child[:len(child)-1]...)
				}

				if throughGroupsOnly(t.Path) {
					newTuples = append(newTuples, t)
				}
			}
		}
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

// e.g. user:alice -> item:1 -> item:2 should not connect user:alice with item:2
// why? because it's not a group
func throughGroupsOnly(path doorman.Path) bool {
	for _, conn := range path {
		if conn.Object.Type() != "group" {
			return false
		}
	}

	return true
}
