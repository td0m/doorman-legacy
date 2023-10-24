package doorman

import (
	"context"
	"strings"
)

type Object string

func (o Object) Type() string {
	return strings.Split(string(o), ":")[0]
}

type ObjectStore interface {
	Add(ctx context.Context, obj Object) error
	Remove(ctx context.Context, obj Object) error
}
