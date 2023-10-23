package server

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
	store "github.com/td0m/doorman/store"
)

var pg *pgxpool.Pool

func cleanup(ctx context.Context) {
	if _, err := pg.Exec(ctx, `delete from tuples`); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}
	pg = pool

	m.Run()
}

func TestCheckStored(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner"},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Store
	tuple := store.Tuple{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	}
	err = db.Add(ctx, tuple)
	require.NoError(t, err)

	// Can access after element was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaRelative(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner"},
					{Label: "viewer", Computed: schema.Relative("owner")},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "viewer",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Store

	tuple := store.Tuple{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	}
	err = db.Add(ctx, tuple)
	require.NoError(t, err)

	// Can access after element was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "viewer",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaRelative2(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "parent"},
					{Label: "owner", Computed: schema.Relative2{From: "parent", Relation: "owner"}},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add

	err = db.Add(ctx, store.Tuple{
		U:     "resource:banana",
		Label: "parent",
		V:     "shop:foo",
	})
	require.NoError(t, err)

	err = db.Add(ctx, store.Tuple{
		U:     "shop:foo",
		Label: "owner",
		V:     "user:alice",
	})
	require.NoError(t, err)

	// Can access after
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaAbsolute(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "group",
				Relations: []schema.Relation{
					{Label: "member"},
				},
			},
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner", Computed: schema.Absolute{U: "group:admins", Label: "member"}},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add

	err = db.Add(ctx, store.Tuple{
		U:     "group:admins",
		Label: "member",
		V:     "user:alice",
	})
	require.NoError(t, err)

	// Can access after
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaAnotherComputed(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner"},
					{Label: "reader", Computed: schema.Relative("owner")},
					{Label: "can_search", Computed: schema.Relative("reader")},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "can_search",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add

	err = db.Add(ctx, store.Tuple{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	require.NoError(t, err)

	// Can access after
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "can_search",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

// TODO: no cyclical computed
