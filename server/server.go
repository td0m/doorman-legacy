package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (dm *Doorman) Connect(ctx context.Context, request *pb.ConnectRequest) (*pb.Relation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (dm *Doorman) Disconnect(ctx context.Context, request *pb.DisconnectRequest) (*pb.Relation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (dm *Doorman) Retrieve(ctx context.Context, request *pb.RetrieveRequest) (*pb.Relation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (dm *Doorman) Check(ctx context.Context, request *pb.CheckRequest) (*pb.CheckResponse, error) {
	fmt.Println("checking", request.Object, request.Name, request.User)

	// First, get stored relations
	relations, err := db.Check(ctx, request.Object, request.Name, request.User)
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

	fromType := extractType(request.Object)
	t, ok := schema.Types[fromType]
	if !ok {
		return nil, fmt.Errorf("invalid type: %s", fromType)
	}

	rel, ok := t.Relations[request.Name]
	if !ok {
		return nil, fmt.Errorf("invalid relation: %s", request.Name)
	}

	success, err := rel.Check(ctx, resolver, request.Object, request.User)
	if err != nil {
		return nil, fmt.Errorf("computing failed: %w", err)
	}

	return &pb.CheckResponse{Connected: success}, nil
}

func extractType(s string) string {
	return strings.SplitN(s, ":", 2)[0]
}
