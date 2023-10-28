package doorman

import (
	"context"
	"fmt"
)

type Relation struct {
	Object  Object
	Verb    Verb
	Subject Object
	Path    []string
}

func (r Relation) String() string {
	return fmt.Sprintf("(%s, %s, %s) via '%s'", r.Subject, r.Verb, r.Object, r.Path)
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
			path := []string{}
			for _, c := range t.Path {
				path = append(path, c.Role, string(c.Object))
			}
			path = append(path, t.Role)
			relations = append(relations, Relation{
				Subject: t.Subject,
				Verb:    verb,
				Object:  t.Object,
				Path:    path,
			})
		}
	}

	return relations, nil
}
