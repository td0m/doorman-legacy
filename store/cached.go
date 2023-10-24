package store

import (
	"context"
	"fmt"

	"github.com/td0m/doorman"
)

type Cached struct {
	store Store

	s2s map[string][]string
	u2s map[doorman.Element][]string
}

func (c Cached) Check(ctx context.Context, set doorman.Set, el doorman.Element) (bool, error) {
	s2s := c.s2s[set.String()]
	u2s := c.u2s[el]
	for _, a := range s2s {
		for _, b := range u2s {
			if a == b {
				fmt.Print("cached!  ")
				return true, nil
			}
		}
	}

	connected, err := c.store.Check(ctx, set, el)
	if err != nil {
		return false, fmt.Errorf("store.Check failed: %w", err)
	}

	return connected, nil
}

func (c Cached) ListElements(ctx context.Context, set doorman.Set) ([]doorman.Element, error) {
	els, err := c.store.ListElements(ctx, set)
	if err != nil {
		return nil, fmt.Errorf("store.ListElements failed: %w", err)
	}

	return els, nil
}

func (c Cached) LogSuccessfulCheck(sets []doorman.Set, v doorman.Element) error {
	for a := len(sets) - 1; a >= 0; a-- {
		for b := a - 1; b >= 0; b-- {
			c.s2s[sets[a].String()] = append(c.s2s[sets[a].String()], sets[b].String())
		}
	}
	for a := len(sets) - 1; a >= 0; a-- {
			c.u2s[v] = append(c.u2s[v], sets[a].String())
	}
	return nil
}

func NewCached(store Store) Cached {
	return Cached{store: store, s2s: map[string][]string{}, u2s: map[doorman.Element][]string{}}
}

func listComputedChangesFromTuple(ctx context.Context, set doorman.Set, v doorman.Element) {
	panic(3)
}
