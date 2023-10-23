package server

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/doorman"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/schema"
	store "github.com/td0m/doorman/store"
)

func newPostgresTuples(ctx context.Context) store.Postgres {
	pool, err := pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}

	return store.NewPostgres(pool)
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestCheckStored(t *testing.T) {
	store.Cleanup()
	ctx := context.Background()

	schema := schema.Schema{}
	server := NewDoormanServer(schema, newPostgresTuples(ctx))

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add edge
	e := edges.Edge{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	}
	err = e.Create(ctx)
	require.NoError(t, err)

	// Can access after edge was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaRelativeEdge(t *testing.T) {
	edges.Cleanup()
	ctx := context.Background()
	server := NewDoormanServer()

	server.schema = doorman.Schema{
		Nodes: map[string]doorman.Node{
			"resource": map[string]doorman.ComputedSet{
				"viewer": doorman.Edge{Node: "", EdgePath: []string{"owner"}},
			},
		},
	}

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "viewer",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add edges
	e1 := edges.Edge{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	}
	err = e1.Create(ctx)
	require.NoError(t, err)

	// Can access after edge was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "viewer",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaRelativeParentEdge(t *testing.T) {
	edges.Cleanup()
	ctx := context.Background()
	server := NewDoormanServer()

	server.schema = doorman.Schema{
		Nodes: map[string]doorman.Node{
			"resource": map[string]doorman.ComputedSet{
				"owner": doorman.Edge{Node: "", EdgePath: []string{"parent", "owner"}},
			},
		},
	}

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add edges
	e2 := edges.Edge{
		U:     "resource:banana",
		Label: "parent",
		V:     "shop:foo",
	}
	err = e2.Create(ctx)
	require.NoError(t, err)

	e1 := edges.Edge{
		U:     "shop:foo",
		Label: "owner",
		V:     "user:alice",
	}
	err = e1.Create(ctx)
	require.NoError(t, err)

	// Can access after edge was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaAbsoluteEdge(t *testing.T) {
	edges.Cleanup()

	ctx := context.Background()
	server := NewDoormanServer()

	server.schema = doorman.Schema{
		Nodes: map[string]doorman.Node{
			"group": map[string]doorman.ComputedSet{
				"member": doorman.Edge{},
			},
			"resource": map[string]doorman.ComputedSet{
				"owner": doorman.Edge{Node: "group:admins", EdgePath: []string{"member"}},
			},
		},
	}

	// Cannot access before
	res, err := server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Add edge
	e := edges.Edge{
		U:     "group:admins",
		Label: "member",
		V:     "user:alice",
	}
	err = e.Create(ctx)
	require.NoError(t, err)

	// Can access after edge was added
	res, err = server.Check(ctx, &pb.CheckRequest{
		U:     "resource:banana",
		Label: "owner",
		V:     "user:alice",
	})
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputedViaAnotherComputedEdge(t *testing.T) {
}
