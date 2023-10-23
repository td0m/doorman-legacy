package schema

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
)

var Nil NilComputed = NilComputed("")

type Schema struct {
	Types []Type
}

func (s Schema) GetRelation(u doorman.Element, relation string) (*Relation, error) {
	// uTypeDef, ok := s.Types[u.Type()]
	// if !ok {
	// 	return nil, fmt.Errorf("failed to get type '%s'", u.Type())
	// }
	//
	// relationDef, ok := uTypeDef[relation]
	// if !ok {
	// 	return nil, fmt.Errorf("failed to relation '%s' for type '%s'", relation, u.Type())
	// }

	uType := u.Type()
	if uType == "" {
		return nil, nil
	}

	for _, t := range s.Types {
		if t.Name == uType {
			for _, r := range t.Relations {
				if r.Label == relation {
					// Do a favour and auto set nil
					if r.Computed == nil {
						r.Computed = Nil
					}
					return &r, nil
				}
			}
			return nil, fmt.Errorf("failed to relation '%s' for type '%s'", relation, u.Type())
		}
	}

	return nil, fmt.Errorf("failed to get type '%s'", u.Type())
}

type Type struct {
	Name      string
	Relations []Relation
}

type Relation struct {
	Label    string
	Computed SetExpr
}

func (r Relation) ToSet(ctx context.Context, rs Resolver, el doorman.Element) (doorman.SetOrOperation, error) {
	return Relative(r.Label).ToSet(ctx, rs, el)
}

type Resolver interface {
	ListElements(ctx context.Context, set doorman.Set) ([]doorman.Element, error)
}

type SetExpr interface {
	ToSet(ctx context.Context, r Resolver, el doorman.Element) (doorman.SetOrOperation, error)
}

type Relative2 struct {
	From     string
	Relation string
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

type Relative string

func (p Relative) ToSet(ctx context.Context, r Resolver, contextualElement doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Set{
		U:     contextualElement,
		Label: string(p),
	}, nil
}

type NilComputed string

func (n NilComputed) ToSet(ctx context.Context, r Resolver, contextualElement doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.NonComputedSet{}, nil
}

type Union struct {
	Exprs []SetExpr
}

func (u Union) ToSet(ctx context.Context, r Resolver, atEl doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Union{}, nil
}

type Absolute doorman.Set

func (s Absolute) ToSet(_ context.Context, _ Resolver, _ doorman.Element) (doorman.SetOrOperation, error) {
	return doorman.Set(s), nil
}
