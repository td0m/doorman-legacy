package doorman

import (
	"context"
	"fmt"
)

type NonComputedSet Set

func (s NonComputedSet) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	directlyContains, err := store.Check(ctx, Set(s), el)
	if err != nil {
		return false, fmt.Errorf("store.Check failed: %w", err)
	}
	return directlyContains, nil
}

type Set struct {
	U     Element
	Label string
}

func (s Set) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	// Reason why we not using a union here rn is performance
	directlyContains, err := store.Check(ctx, s, el)
	if err != nil {
		return false, fmt.Errorf("store.Check failed: %w", err)
	}
	if directlyContains {
		return true, nil
	}

	computed, err := store.Computed(ctx, s)
	if err != nil {
		return false, fmt.Errorf("store.ListSubsets failed: %w", err)
	}

	if computed == nil {
		return false, nil
	}

	return computed.Contains(ctx, store, el)
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

type Intersection []SetOrOperation

func (i Intersection) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	for _, setOrOp := range i {
		contains, err := setOrOp.Contains(ctx, store, el)
		if err != nil {
			return false, err
		}
		if !contains {
			return false, nil
		}
	}
	return true, nil
}
