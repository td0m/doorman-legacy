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

// no cycles, valid references
// TODO: func (s Schema) Validate() error {}

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

type Resolver interface {
	ListElements(ctx context.Context, set doorman.Set) ([]doorman.Element, error)
}
