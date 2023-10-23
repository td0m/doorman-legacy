package store

import (
	"context"

	"github.com/td0m/doorman"
)

type Store interface {
	Check(ctx context.Context, s doorman.Set, e doorman.Element) (bool, error)
	ListElements(ctx context.Context, set doorman.Set) ([]doorman.Element, error)
}
