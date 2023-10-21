package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
)

type dbResolver struct {
}

func (r dbResolver) ListForward(ctx context.Context, from, name string) ([]db.Relation, error) {
	return db.ListForward(ctx, from, name)
}

var resolver = dbResolver{}

func NewDoormanServer() *Doorman {
	return &Doorman{}
}

type Doorman struct {
	*pb.UnimplementedDoormanServer
}

// TODO: make it not idempotent? no two rels with same name? needed for good caching
func (dm *Doorman) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.Relation, error) {
	r := &db.Relation{
		From: request.From,
		Name: request.Name,
		To:   request.To,
	}
	if err := r.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed: %w", err)
	}
	return &pb.Relation{
		From: r.From,
		Name: r.Name,
		To:   r.To,
	}, nil
}

func (dm *Doorman) Disconnect(ctx context.Context, request *pb.DisconnectRequest) (*pb.Relation, error) {
	r := &db.Relation{
		From: request.From,
		Name: request.Name,
		To:   request.To,
	}
	// if err := r.Delete(ctx); err != nil {
	// 	return nil, fmt.Errorf("db failed: %w", err)
	// }
	return &pb.Relation{
		From: r.From,
		Name: r.Name,
		To:   r.To,
	}, nil
}

func (dm *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	fmt.Println("checking", request.From, request.Name, request.To)

	// First, get stored relations
	relations, err := db.Check(ctx, request.From, request.Name, request.To)
	if err != nil {
		return nil, fmt.Errorf("db failed: %w", err)
	}

	if len(relations) > 0 {
		return &pb.CheckResponse{Connected: true}, nil
	}

	// Obtain computed relations next (if not successful?)
	schema := schema.Schema{
		Types: map[string]schema.Type{
			"product": {
				Relations: map[string]schema.SetExpr{
					"owner": schema.Union{[]schema.SetExpr{schema.Path([]string{"", "on", "owner"}), schema.Path([]string{"", "foo"})}},
				},
			},
		},
	}

	fromType := extractType(request.From)
	t, ok := schema.Types[fromType]
	if !ok {
		return nil, fmt.Errorf("invalid type: %s", fromType)
	}

	rel, ok := t.Relations[request.Name]
	if !ok {
		return nil, fmt.Errorf("invalid relation: %s", request.Name)
	}

	success, err := rel.Check(ctx, resolver, request.From, request.To)
	if err != nil {
		return nil, fmt.Errorf("computing failed: %w", err)
	}

	return &pb.CheckResponse{Connected: success}, nil
}

func extractType(s string) string {
	return strings.SplitN(s, ":", 2)[0]
}
