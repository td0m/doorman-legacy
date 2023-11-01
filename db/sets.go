package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/td0m/doorman"
)

type Sets struct {
	conn            querier
	subject2parents map[doorman.Object]sets
	// recursive subsets
	set2subset   map[doorman.Set]sets
	set2superset map[doorman.Set]sets
}

type sets struct {
	m     map[doorman.Set]bool
	stale bool
}

func newSets() sets {
	return sets{
		m:     map[doorman.Set]bool{},
		stale: false,
	}
}

func (s sets) Add(set doorman.Set) {
	s.m[set] = true
}

func (s sets) Remove(set doorman.Set) {
	s.m[set] = false
}

var ErrStale = errors.New("stale cache")

func (s Sets) Contains(ctx context.Context, set doorman.Set, subject doorman.Object) (bool, error) {
	parents, ok := s.subject2parents[subject]
	if !ok || parents.stale {
		return false, ErrStale
	}

	subsets, ok := s.set2subset[set]
	if !ok || subsets.stale {
		// CONSIDER
		self := newSets()
		self.Add(set)
		if intersect(parents, self) {
			fmt.Println("cache!")
			return true, nil
		} else {
			return false, ErrStale
		}
	}

	return intersect(parents, subsets), nil
}

func (s Sets) UpdateParents(ctx context.Context, subject doorman.Object, sets []doorman.Set) error {
	s.subject2parents[subject] = setsFromList(sets)
	return nil
}

func (s Sets) InvalidateParents(ctx context.Context, subject doorman.Object) error {
	sets := newSets()
	sets.stale = true
	s.subject2parents[subject] = sets
	return nil
}


func (s Sets) UpdateSubsets(ctx context.Context, set, subsets []doorman.Set) error {
	// return s.modifySubset(ctx, set, subset, true)
	return nil
}

func setsFromList(ss []doorman.Set) sets {
	sets := newSets()
	for _, s := range ss {
		sets.Add(s)
	}
	return sets
}

func intersect(a, b sets) bool {
	for a, aconn := range a.m {
		if aconn && b.m[a] {
			return true
		}
	}
	return false
}

func NewSets(q querier) Sets {
	return Sets{
		conn:            q,
		subject2parents: map[doorman.Object]sets{},
		set2subset:      map[doorman.Set]sets{},
		set2superset:    map[doorman.Set]sets{},
	}
}
