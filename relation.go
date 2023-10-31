package doorman

import (
	"context"
	"fmt"
)

type Relation struct {
	Object  Object `json:"object"`
	Verb    Verb   `json:"verb"`
	Subject Object `json:"subject"`
}

func (r Relation) String() string {
	return fmt.Sprintf("(%s, %s, %s)", r.Subject, r.Verb, r.Object)
}

type resolveRole func(ctx context.Context, id string) (*Role, error)

func TuplesToRelations(ctx context.Context, tuples []Tuple, r resolveRole) ([]Relation, error) {
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
		if !t.Path.GroupsOnly() {
			continue
		}
		for _, verb := range uniqueRoles[t.Role].Verbs {
			relations = append(relations, Relation{
				Subject: t.Subject,
				Verb:    verb,
				Object:  t.Object,
			})
		}
	}

	return relations, nil
}
