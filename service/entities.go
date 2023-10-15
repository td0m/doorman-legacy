package service

import (
	"context"
	"fmt"

	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Entities struct {
	*pb.UnimplementedEntitiesServer
}

func (es *Entities) Create(ctx context.Context, request *pb.EntitiesCreateRequest) (*pb.Entity, error) {
	e := &db.Entity{
		ID:    request.Id,
		Attrs: request.Attrs.AsMap(),
	}
	if err := e.Create(ctx); err != nil {
		return nil, fmt.Errorf("db.Create failed: %w", err)
	}

	return mapEntityFromDB(*e), nil
}

func (es *Entities) Retrieve(ctx context.Context, request *pb.EntitiesRetrieveRequest) (*pb.Entity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}

func (es *Entities) List(ctx context.Context, request *pb.EntitiesListRequest) (*pb.EntitiesListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func (es *Entities) Update(ctx context.Context, request *pb.EntitiesUpdateRequest) (*pb.Entity, error) {
	entity, err := db.RetrieveEntity(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("db.RetrieveEntity failed: %w", err)
	}

	// TODO: upsert
	if request.Attrs != nil {
		entity.Attrs = request.Attrs.AsMap()
	}

	if err := entity.Update(ctx); err != nil {
		return nil, fmt.Errorf("db.Update failed: %w", err)
	}

	return mapEntityFromDB(*entity), nil
}

func (es *Entities) Delete(ctx context.Context, request *pb.EntitiesDeleteRequest) (*pb.EntitiesDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func mapEntityFromDB(e db.Entity) *pb.Entity {
	attrs, _ := structpb.NewStruct(e.Attrs)
	return &pb.Entity{
		Id: e.ID,
		Attrs: attrs,
	}
}
