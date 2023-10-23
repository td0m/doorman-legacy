package server

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewDoormanServer(sch schema.Schema, store doorman.Store) *Doorman {
	return &Doorman{
		Schema: sch,
		Store:  store,
	}
}

type Doorman struct {
	*pb.UnimplementedDoormanServer
	Schema schema.Schema
	Store  doorman.Store
}

func (d *Doorman) Write(ctx context.Context, request *pb.WriteRequest) (*pb.WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}

func (d *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	u := doorman.Element(request.U)
	v := doorman.Element(request.V)

	relationDef, err := d.Schema.GetRelation(u, request.Label)
	if err != nil {
		return nil, fmt.Errorf("schema failed to get relation: %w", err)
	}

	set, err := relationDef.ToSet(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("schema relationDef.ToSet failed: %w", err)
	}

	contains, err := set.Contains(ctx, d.Store, v)
	if err != nil {
		return nil, fmt.Errorf("computedSet.Contains failed: %w", err)
	}

	return &pb.CheckResponse{
		Connected: contains,
	}, nil
}

