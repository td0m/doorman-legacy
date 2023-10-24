package doorman

import "context"

type Tuple struct {
	Subject Object
	Role    string
	Object  Object
}

func NewTuple(sub Object, role string, obj Object) Tuple {
	return Tuple{sub, role, obj}
}

type Connection struct {
	Role   string
	Object Object
}

type Path []Connection

type TupleStore interface {
	Add(ctx context.Context, t Tuple) error
	Remove(ctx context.Context, t Tuple) error
	ListConnected(ctx context.Context, subject Object, inverted bool) ([]Path, error)
}
