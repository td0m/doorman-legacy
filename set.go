package doorman

import (
	"context"
	"fmt"
)

type Exclusion struct {
	A SetOrOperation
	B SetOrOperation
}

func (e Exclusion) Contains(ctx context.Context, store Store, el Element) (bool, []Set, error) {
	Acontains, path, err := e.A.Contains(ctx, store, el)
	if err != nil {
		return false, nil, fmt.Errorf("A.Contains failed: %w", err)
	}
	if !Acontains {
		return false, nil, nil
	}

	Bcontains, _, err := e.B.Contains(ctx, store, el)
	if err != nil {
		return false, nil, fmt.Errorf("B.Contains failed: %w", err)
	}

	if Bcontains {
		return false, nil, nil
	}

	return true, path, nil
}

type Intersection []SetOrOperation

func (i Intersection) Contains(ctx context.Context, store Store, el Element) (bool, []Set, error) {
	var lastPath []Set
	for _, setOrOp := range i {
		contains, path, err := setOrOp.Contains(ctx, store, el)
		if err != nil {
			return false, nil, err
		}
		if !contains {
			return false, nil, nil
		}
		lastPath = path
	}
	// TODO: consider how to log all paths...
	return true, lastPath, nil
}

type NonComputedSet Set

func (s NonComputedSet) Contains(ctx context.Context, store Store, el Element) (bool, []Set, error) {
	directlyContains, err := store.Check(ctx, Set(s), el)
	if err != nil {
		return false, nil, fmt.Errorf("store.Check failed: %w", err)
	}
	return directlyContains, []Set{Set(s)}, nil
}

type Set struct {
	U     Element
	Label string
}

func (s Set) String() string {
	return string(s.U) + "." + string(s.Label)
}

func (s Set) Contains(ctx context.Context, store Store, el Element) (bool, []Set, error) {
	// Reason why we not using a union here rn is performance
	directlyContains, err := store.Check(ctx, s, el)
	if err != nil {
		return false, nil, fmt.Errorf("store.Check failed: %w", err)
	}
	if directlyContains {
		return true, []Set{Set(s)}, nil
	}

	computed, err := store.Computed(ctx, s)
	if err != nil {
		return false, nil, fmt.Errorf("store.ListSubsets failed: %w", err)
	}

	if computed == nil {
		return false, nil, nil
	}

	computedContains, path, err := computed.Contains(ctx, store, el)
	if err != nil {
		return false, nil, fmt.Errorf("computed.Contains failed: %w", err)
	}

	if !computedContains {
		return false, nil, nil
	}

	return true, append(path, Set(s)), nil
}

type SetOrOperation interface {
	Contains(ctx context.Context, store Store, el Element) (bool, []Set, error)
}

type Union []SetOrOperation

func (u Union) Contains(ctx context.Context, store Store, el Element) (bool, []Set, error) {
	for _, setOrOp := range u {
		contains, path, err := setOrOp.Contains(ctx, store, el)
		if err != nil {
			return false, nil, err
		}
		if contains {
			return contains, path, nil
		}
	}
	return false, nil, nil
}
