package doorman

import "context"

type Verb string

type Role struct {
	ID    string
	Verbs []Verb
}

func NewRole(id string, optverbs ...[]Verb) Role {
	verbs := []Verb{}
	if len(optverbs) > 0 {
		verbs = optverbs[0]
	}
	return Role{ID: id, Verbs: verbs}
}

type RoleStore interface {
	Add(ctx context.Context, role Role) error
	Remove(ctx context.Context, id string) error
	Retrieve(ctx context.Context, id string) (*Role, error)
	UpdateOne(ctx context.Context, role Role) error
}
