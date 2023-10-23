package server

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
	"github.com/td0m/doorman/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type setStore struct {
	schema schema.Schema
	tuples store.Postgres
}

func (s setStore) ListSubsets(ctx context.Context, set doorman.Set) (doorman.SetOrOperation, error) {
	rel, err := s.schema.GetRelation(set.U, set.Label)
	if err != nil {
		return nil, fmt.Errorf("GetRelation failed: %w", err)
	}
	return rel.Computed.ToSet(ctx, set.U)
}

func (s setStore) Check(ctx context.Context, set doorman.Set, el doorman.Element) (bool, error) {
	return s.tuples.Check(ctx, set, el)
}

func NewDoormanServer(sch schema.Schema, store store.Postgres) *Doorman {
	return &Doorman{
		setStore: setStore{schema: sch, tuples: store},
	}
}

type Doorman struct {
	*pb.UnimplementedDoormanServer
	setStore setStore
}

func (d *Doorman) Write(ctx context.Context, request *pb.WriteRequest) (*pb.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}

func (d *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	u := doorman.Element(request.U)
	v := doorman.Element(request.V)

	relationDef, err := d.setStore.schema.GetRelation(u, request.Label)
	if err != nil {
		return nil, fmt.Errorf("schema failed to get relation: %w", err)
	}

	set, err := relationDef.ToSet(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("schema relationDef.ToSet failed: %w", err)
	}

	contains, err := set.Contains(ctx, d.setStore, v)
	if err != nil {
		return nil, fmt.Errorf("computedSet.Contains failed: %w", err)
	}

	return &pb.CheckResponse{
		Connected: contains,
	}, nil
}

