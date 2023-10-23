package schema

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
)

var Nil NilComputed = NilComputed{}

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

	for _, t := range s.Types {
		if t.Name == u.Type() {
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

func (r Relation) ToSet(ctx context.Context, el doorman.Element) (doorman.Set, error) {
	return RelativePath{r.Label}.ToSet(ctx, el)
}

type SetExpr interface {
	ToSet(ctx context.Context, el doorman.Element) (doorman.Set, error)
}

type RelativePath []string

func (p RelativePath) ToSet(ctx context.Context, contextualElement doorman.Element) (doorman.Set, error) {
	return doorman.Set{
		U:     contextualElement,
		Label: p[0],
	}, nil
}

type NilComputed []string

func (NilComputed) ToSet(ctx context.Context, contextualElement doorman.Element) (doorman.Set, error) {
	return doorman.Set{}, nil
}

type Union struct {
	Exprs []SetExpr
}

func (u Union) ToSet(ctx context.Context, atEl doorman.Element) (doorman.Union, error) {
	return doorman.Union{}, nil
}
