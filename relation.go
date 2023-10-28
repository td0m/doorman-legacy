package doorman

import (
	"context"
	"fmt"
	"strings"
)

type Relation struct {
	Object  Object
	Verb    Verb
	Subject Object
	Key     string
}

func (r Relation) String() string {
	return fmt.Sprintf("(%s, %s, %s) via '%s'", r.Subject, r.Verb, r.Object, r.Key)
}

type resolveRole func(ctx context.Context, id string) (*Role, error)

func TuplesToRelations(ctx context.Context, tuples []TupleWithPath, r resolveRole) ([]Relation, error) {
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

	relations := []Relation{}
	for _, t := range tuples {
		for _, verb := range uniqueRoles[t.Role].Verbs {
			s := []string{}
			for _, c := range t.Path {
				s = append(s, ">"+c.Role+">"+string(c.Object))
			}
			relations = append(relations, Relation{
				Subject: t.Subject,
				Verb:    verb,
				Object:  t.Object,
				Key:     strings.Join(s, ", "),
			})
		}
	}

	return relations, nil
}
