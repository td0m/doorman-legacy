package doorman

import (
	"context"
	"fmt"
)

type Set struct {
	Object Object
	Verb   Verb
}

func (s Set) String() string {
	return fmt.Sprintf("%s.%s", s.Object, s.Verb)
}

func NewSet(o Object, verb Verb) Set {
	return Set{Object: o, Verb: verb}
}

type resolveRole func(ctx context.Context, id string) (*Role, error)

func ParentTuplesToSets(ctx context.Context, tuples []Tuple, r resolveRole) ([]Set, error) {
	uniqueRoles := map[string]*Role{}
	for _, t := range tuples {
		uniqueRoles[t.Role] = nil
	}

	for id := range uniqueRoles {
		role, err := r(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("resolving role failed: %w", err)
		}
		uniqueRoles[id] = role
	}

	sets := []Set{}
	for _, t := range tuples {
		for _, verb := range uniqueRoles[t.Role].Verbs {
			sets = append(sets, Set{
				Object: t.Object,
				Verb:   verb,
			})
		}
	}

	return sets, nil
}
