package schema

import (
	"context"
	"github.com/td0m/doorman/db"
)


type Schema struct {
	Types map[string]Type
}

type Type struct {
	Relations map[string]SetExpr
}

type SetExpr interface {
	Check(ctx context.Context, r Resolver, from, to string) (bool, error)
}

type Resolver interface {
	ListForward(ctx context.Context, from, name string) ([]db.Relation, error)
}

