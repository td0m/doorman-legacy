package doorman

import (
	"context"
)

type Store interface {
	List(ctx context.Context, set Set) ([]string, error)
	Check(ctx context.Context, set Set, element Element) (bool, error)
}

