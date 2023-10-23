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

func TestCheckComputedUnion(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner"},
					{Label: "manager"},
					{Label: "viewer", Computed: schema.Union{schema.Relative("owner"), schema.Relative("manager")}},
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

	t.Run("True if both connect", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{
			U:     "resource:banana",
			Label: "owner",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		tupleB := store.Tuple{
			U:     "resource:banana",
			Label: "manager",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleB)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "viewer",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Connected)
	})

	t.Run("True if only one connects", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{
			U:     "resource:banana",
			Label: "owner",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "viewer",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Connected)
	})
}

func TestCheckComputedIntersection(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "owner"},
					{Label: "manager"},
					{Label: "viewer", Computed: schema.Intersection{schema.Relative("owner"), schema.Relative("manager")}},
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

	t.Run("True if both connect", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{
			U:     "resource:banana",
			Label: "owner",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		tupleB := store.Tuple{
			U:     "resource:banana",
			Label: "manager",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleB)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "viewer",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Connected)
	})

	t.Run("False if only one connects", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{
			U:     "resource:banana",
			Label: "owner",
			V:     "user:alice",
		}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "viewer",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Connected)
	})
}

func TestCheckComputedExclusion(t *testing.T) {
	ctx := context.Background()
	cleanup(ctx)

	schema := schema.Schema{
		Types: []schema.Type{
			{
				Name: "resource",
				Relations: []schema.Relation{
					{Label: "a"},
					{Label: "b"},
					{Label: "c", Computed: schema.Exclusion{A: schema.Relative("a"), B: schema.Relative("b")}},
				},
			},
		},
	}
	db := store.NewPostgres(pg)
	server := NewDoormanServer(schema, db)

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "c",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	t.Run("False if both connect", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{U: "resource:banana", Label: "a", V: "user:alice"}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		tupleB := store.Tuple{U: "resource:banana", Label: "b", V: "user:alice"}
		err = db.Add(ctx, tupleB)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "c",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Connected)
	})

	t.Run("True if only A connects", func(t *testing.T) {
		cleanup(ctx)

		tupleA := store.Tuple{U: "resource:banana", Label: "a", V: "user:alice"}
		err = db.Add(ctx, tupleA)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "c",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Connected)
	})

	t.Run("False if only B connects", func(t *testing.T) {
		cleanup(ctx)

		tupleB := store.Tuple{U: "resource:banana", Label: "b", V: "user:alice"}
		err = db.Add(ctx, tupleB)
		require.NoError(t, err)

		res, err = server.Check(ctx, &pb.CheckRequest{
			U:     "resource:banana",
			Label: "c",
			V:     "user:alice",
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Connected)
	})
}

// TODO: no cyclical computed
