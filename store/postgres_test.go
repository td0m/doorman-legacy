package tuples

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}
	pg = pool

	m.Run()
}

func TestTupleCreateSuccess(t *testing.T) {
	Cleanup()
	ctx := context.Background()

	e := &Tuple{
		U: "resource:u",
		Label: "owner",
		V: "user:alice",
	}

	err := e.Create(ctx)
	assert.NoError(t, err)
}

