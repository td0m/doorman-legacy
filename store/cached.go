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
	fmt.Println("can i cache?")

	connected, err := c.store.Check(ctx, set, el)
	if err != nil {
		return false, fmt.Errorf("store.Check failed: %w", err)
	}

	if connected {
		c.u2s[el] = append(c.u2s[el], set.String())
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

func NewCached(store Store) Cached {
	return Cached{store: store, s2s: map[string][]string{}, u2s: map[doorman.Element][]string{} }
}
