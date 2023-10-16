package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validConnections = [][]string{}
var registeredTypes = []string{}

func Setup(confValidConnections [][]string) error {
	registeredTypes = []string{}
	validConnections = confValidConnections
	for _, vs := range validConnections {
		if len(vs) != 2 {
			return errors.New("expected each connection to contain exactly two items")
		}
		registeredTypes = append(registeredTypes, vs...)
	}

	return nil
}

type Relations struct {
	*pb.UnimplementedRelationsServer
}

func (rs *Relations) Create(ctx context.Context, request *pb.RelationsCreateRequest) (*pb.Relation, error) {
	fromType, toType := extractType(request.FromId), extractType(request.ToId)
	if !isValidConnection(fromType, toType) {
		return nil, status.Errorf(codes.InvalidArgument, "relation between types '%s' and '%s' is not allowed", fromType, toType)
	}
	r := &db.Relation{
		From: request.FromId,
		To:   request.ToId,
		Name: request.Name,
	}
	if err := r.Create(ctx); err != nil {
		return nil, fmt.Errorf("create failed: %w", err)
	}

	return mapRelationFromDB(*r), nil
}
func (rs *Relations) Retrieve(ctx context.Context, request *pb.RelationsRetrieveRequest) (*pb.Relation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (rs *Relations) List(ctx context.Context, request *pb.RelationsListRequest) (*pb.RelationsListResponse, error) {
	table := "cache"
	if request.NoCache != nil && !*request.NoCache {
		table = "relations"
	}

	// TODO: if from or to specified we validate it here

	f := db.RelationFilter{
		AfterID:  request.PaginationToken,
		From:     request.FromId,
		FromType: request.FromType,
		Name:     request.Name,
		To:       request.ToId,
		ToType:   request.ToType,
	}

	relations, err := db.ListRelationsOrCache(ctx, table, f)
	if err != nil {
		return &pb.RelationsListResponse{}, fmt.Errorf("db failed: %w", err)
	}

	items := make([]*pb.Relation, len(relations))
	for i := range relations {
		items[i] = mapRelationFromDBCache(relations[i])
	}

	return &pb.RelationsListResponse{
		Items: items,
	}, nil
}
func (rs *Relations) Update(ctx context.Context, request *pb.RelationsUpdateRequest) (*pb.Relation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (rs *Relations) Delete(ctx context.Context, request *pb.RelationsDeleteRequest) (*pb.RelationsDeleteResponse, error) {
	r, err := db.RetrieveRelation(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("db.RetrieveRelation failed: %w", err)
	}

	if err := r.Delete(ctx); err != nil {
		return nil, fmt.Errorf("reaction.Delete failed: %w", err)
	}
	return &pb.RelationsDeleteResponse{}, nil
}

func mapRelationFromDB(r db.Relation) *pb.Relation {
	return &pb.Relation{
		Id:     r.ID,
		FromId: r.From,
		ToId:   r.To,
		Name:   r.Name,
	}
}

func mapRelationFromDBCache(r db.Cache) *pb.Relation {
	return &pb.Relation{
		Id:     r.ID,
		FromId: r.From,
		ToId:   r.To,
		Name:   r.Name,
	}
}

func extractType(id string) string {
	return strings.Split(id, ":")[0]
}

func isValidConnection(from, to string) bool {
	if len(validConnections) == 0 {
		return true
	}

	for _, c := range validConnections {
		if c[0] == from {
			if c[1] == to {
				return true
			} else if c[1] == "*" && !slices.Contains(registeredTypes, to) {
				return true
			}
		}
	}

	return false
}
