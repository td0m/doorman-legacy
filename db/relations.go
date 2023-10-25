package db

import (
	"context"

	"github.com/td0m/doorman"
)

type Relations struct {
	// TODO: define a group?
	group2group   map[group]SetOfGroups
	subject2group map[doorman.Object]SetOfGroups
}

func (rs Relations) Check(ctx context.Context, r doorman.Relation) (bool, error) {
	subject2group, ok := rs.subject2group[r.Subject]
	if !ok {
		return false, nil
	}
	group2group, ok := rs.group2group[group{Object: r.Object, Verb: r.Verb}]
	if !ok {
		return false, nil
	}

	return subject2group.Intersects(group2group), nil
}

func (rs Relations) Add(ctx context.Context, r doorman.Relation) error {
	g := group{Object: r.Object, Verb: r.Verb}

	subject2group, ok := rs.subject2group[r.Subject]
	if !ok {
		subject2group = SetOfGroups{}
	}
	subject2group[g] = true

	group2group, ok := rs.group2group[g]
	if !ok {
		group2group = SetOfGroups{}
	}
	group2group[g] = true // connnect to self

	// TODO: also take care of g2g

	// Save
	rs.subject2group[r.Subject] = subject2group
	rs.group2group[g] = group2group

	return nil
}

func (rs Relations) Remove(ctx context.Context, r doorman.Relation) error {
	panic("fak")
}

func NewRelations() Relations {
	return Relations{
		group2group:   map[group]SetOfGroups{},
		subject2group: map[doorman.Object]SetOfGroups{},
	}
}

type group struct {
	Object doorman.Object
	Verb   doorman.Verb
}

type SetOfGroups map[group]bool

func (a SetOfGroups) Intersects(b SetOfGroups) bool {
	for b, bconn := range b {
		if bconn && a[b] {
			return true
		}
	}
	return false
}
