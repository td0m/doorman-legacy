package schema

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
)

type Absolute doorman.Set

type Exclusion struct {
	A SetExpr
	B SetExpr
}

type Intersection []SetExpr

type NilComputed string

type Relation struct {
	Label    string
	Computed SetExpr
}

type Relative string

type Relative2 struct {
	From     string
	Relation string
}

type SetExpr interface {
	ToSet(ctx context.Context, r Resolver, el doorman.Element) (doorman.SetOrOperation, error)
}

type Union []SetExpr

func (s Absolute) ToSet(_ context.Context, _ Resolver, _ doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Set(s), nil
}

func (e Exclusion) ToSet(ctx context.Context, r Resolver, atEl doorman.Element) (doorman.SetOrOperation, error) {
	setA, err := e.A.ToSet(ctx, r, atEl)
	if err != nil {
		return nil, fmt.Errorf("set A failed: %w", err)
	}

	setB, err := e.B.ToSet(ctx, r, atEl)
	if err != nil {
		return nil, fmt.Errorf("set B failed: %w", err)
	}
	return doorman.Exclusion{A: setA, B: setB}, nil
}

func (i Intersection) ToSet(ctx context.Context, r Resolver, atEl doorman.Element) (doorman.SetOrOperation, error) {
	sets := make([]doorman.SetOrOperation, len(i))
	for i, v := range i {
		set, err := v.ToSet(ctx, r, atEl)
		if err != nil {
			return nil, fmt.Errorf("child %d failed: %w", i, err)
		}
		sets[i] = set
	}
	return doorman.Intersection(sets), nil
}

func (n NilComputed) ToSet(ctx context.Context, r Resolver, contextualElement doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.NonComputedSet{}, nil
}

func (r Relation) ToSet(ctx context.Context, rs Resolver, el doorman.Element) (doorman.SetOrOperation, error) {
	return Relative(r.Label).ToSet(ctx, rs, el)
}

func (p Relative) ToSet(ctx context.Context, _ Resolver, contextualElement doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Set{
		U:     contextualElement,
		Label: string(p),
	}, nil
}

func (p Relative2) ToSet(ctx context.Context, r Resolver, contextualElement doorman.Element) (doorman.SetOrOperation, error) {
	relations, err := r.ListElements(ctx, doorman.Set{U: contextualElement, Label: p.From})
	if err != nil {
		return nil, err
	}

	sets := make([]doorman.SetOrOperation, len(relations))
	for i, r := range relations {
		sets[i] = doorman.Set{U: r, Label: p.Relation}
	}

	return doorman.Union(sets), nil
}

func (u Union) ToSet(ctx context.Context, r Resolver, atEl doorman.Element) (doorman.SetOrOperation, error) {
	sets := make([]doorman.SetOrOperation, len(u))
	for i, v := range u {
		set, err := v.ToSet(ctx, r, atEl)
		if err != nil {
			return nil, fmt.Errorf("child %d failed: %w", i, err)
		}
		sets[i] = set
	}
	return doorman.Union(sets), nil
}
