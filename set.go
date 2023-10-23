package doorman

import (
	"context"
	"fmt"
)

type Set struct {
	U     Element
	Label string
}

func (s Set) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	directlyContains, err := store.Check(ctx, s, el)
	if err != nil {
		return false, fmt.Errorf("store.Check failed: %w", err)
	}
	if directlyContains {
		return true, nil
	}
	return false, nil
}

type SetOrOperation interface {
	Contains(ctx context.Context, store Store, el Element) (bool, error)
}

type Union []SetOrOperation

func (u Union) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	for _, setOrOp := range u {
		contains, err := setOrOp.Contains(ctx, store, el)
		if err != nil {
			return false, err
		}
		if contains {
			return contains, nil
		}
	}
	return false, nil
}
