package doorman

import (
	"context"
)

type Set struct {
	U string
	Label string
}

func (s Set) Contains(ctx context.Context, store Store, el Element) (bool, error) {
	return false, nil
}
