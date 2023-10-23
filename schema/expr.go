package schema

import (
	"context"

	"github.com/td0m/doorman"
)

type Absolute doorman.Set

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

type Union struct {
	Exprs []SetExpr
}

func (s Absolute) ToSet(_ context.Context, _ Resolver, _ doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Set(s), nil
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
	return doorman.Union{}, nil
}
