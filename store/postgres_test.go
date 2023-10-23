package store

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/td0m/doorman"
)

var pg *pgxpool.Pool

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
	db := NewPostgres(pg)

	ctx := context.Background()

	s := doorman.Set{
		U:     doorman.Element("resource:rnd" + xid.New().String()),
		Label: "owner",
	}

	err := db.Add(ctx, s, doorman.Element("user:alice"))
	assert.NoError(t, err)
}
