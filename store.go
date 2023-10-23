package doorman

import (
	"context"
)

type SetEnumerator interface {
	List(ctx context.Context, set Set) ([]Element, error)
}

type SetChecker interface {
	Check(ctx context.Context, set Set, element Element) (bool, error)
}

type SetModifier interface {
	Add(ctx context.Context, set Set, element Element) error
	// Remove(ctx context.Context, set Set, element Element) error
}

type Store interface {
	SetModifier
	SetChecker
	SetEnumerator
}
