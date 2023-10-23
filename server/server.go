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
	tuples store.Store
}

func (s setStore) Computed(ctx context.Context, set doorman.Set) (doorman.SetOrOperation, error) {
	rel, err := s.schema.GetRelation(set.U, set.Label)
	if err != nil {
		return nil, fmt.Errorf("GetRelation failed: %w", err)
	}
	return rel.Computed.ToSet(ctx, s.tuples, set.U)
}

func (s setStore) Check(ctx context.Context, set doorman.Set, el doorman.Element) (bool, error) {
	return s.tuples.Check(ctx, set, el)
}

func NewDoormanServer(sch schema.Schema, store store.Store) *Doorman {
	return &Doorman{
		setStore: setStore{schema: sch, tuples: store},
	}
}

type Doorman struct {
	*pb.UnimplementedDoormanServer
	setStore setStore

	cache *store.Cached
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

	set, err := relationDef.ToSet(ctx, d.setStore.tuples, u)
	if err != nil {
		return nil, fmt.Errorf("schema relationDef.ToSet failed: %w", err)
	}

	contains, path, err := set.Contains(ctx, d.setStore, v)
	if err != nil {
		return nil, fmt.Errorf("computedSet.Contains failed: %w", err)
	}

	if d.cache != nil && contains {
		_ = d.cache.LogSuccessfulCheck(path, v)
	}

	return &pb.CheckResponse{
		Connected: contains,
	}, nil
}
