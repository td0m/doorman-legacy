package schema

import (
	"context"
	"fmt"
)

type Union struct {
	Expr []SetExpr
}

func (u Union) Check(ctx context.Context, r Resolver, from, to string) (bool, error) {
	// TODO: caching to improve performance
	for _, expr := range u.Expr {
		succ, err := expr.Check(ctx, r, from, to)
		if err != nil {
			return false, fmt.Errorf("resolve failed: %w", err)
		}
		if succ {
			return true, nil
		}
	}
	return false, nil
}

