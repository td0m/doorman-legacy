package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	schema := Schema{
		Types: map[string]Type{
			"product": {
				Relations: map[string]SetExpr{
					"owner": Union{[]SetExpr{Path([]string{"", "on", "owner"}), Path([]string{"", "foo"})}},
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

	success, err := rel.Check(ctx, request.Object)
	if err != nil {
		return nil, fmt.Errorf("computing failed: %w", err)
	}

	return &pb.CheckResponse{Connected: success}, nil
}

// TODO: move these types to "schema" package, store schema in memory / db

type Schema struct {
	Types map[string]Type
}

type Type struct {
	Relations map[string]SetExpr
}

type SetExpr interface {
	Check(ctx context.Context, from string) (bool, error)
}

type Union struct {
	Expr []SetExpr
}

func (u Union) Check(ctx context.Context, from string) (bool, error) {
	for _, expr := range u.Expr {
		succ, err := expr.Check(ctx, from)
		if err != nil {
			return false, fmt.Errorf("resolve failed: %w", err)
		}
		if succ {
			return true, nil
		}
	}
	return false, nil
}

type Path []string

func (path Path) Check(ctx context.Context, from string) (bool, error) {
	if len(path) < 2 {
		return false, fmt.Errorf("bad path length")
	}

	if len(path) > 0 && path[0] == "" {
		path[0] = from
	}

	relations, err := db.ListForward(ctx, path[0], path[1])
	if err != nil {
		return false, fmt.Errorf("ListForward failed: %w", err)
	}

	// TODO: consider if this can be made concurrent. I think yes.
	for _, relation := range relations {
		subpath := path[1:]
		if len(subpath) == 1 {
			// REACHED END OF PATH! JUST RETURN
			return true, nil
		}
		subpath[0] = ""

		succ, err := Path(subpath).Check(ctx, relation.To)
		if err != nil {
			return false, fmt.Errorf("path resolve failed: %w", err)
		}

		// Full path at any point means can return true
		if succ {
			return true, nil
		}
	}
	return false, nil
}

// TODO: intersection, negation

func extractType(s string) string {
	return strings.SplitN(s, ":", 2)[0]
}
