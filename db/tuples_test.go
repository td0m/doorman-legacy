package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/doorman"
)

var conn *pgxpool.Pool

func cleanup(conn *pgxpool.Pool) {
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		delete from tuples;
		delete from roles;
		delete from objects;
	`)

	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "")
	if err != nil {
		panic(err)
	}
	conn = pool

	m.Run()
}

func TestTuplesCreate(t *testing.T) {
	ctx := context.Background()
	cleanup(conn)
	tuples := NewTuples(conn)
	objects := NewObjects(conn)
	roles := NewRoles(conn)

	// Make objects and roles

	manisha := doorman.Object("user:manisha")
	require.NoError(t, objects.Add(ctx, manisha))

	member := doorman.Role{
		ID:    "member",
		Verbs: []doorman.Verb{"foo"},
	}
	require.NoError(t, roles.Add(ctx, member))

	admins := doorman.Object("group:admins")
	require.NoError(t, objects.Add(ctx, admins))

	// Create tuple

	t.Run("Success creating tuple between existing items", func(t *testing.T) {
		tuple := doorman.NewTuple(manisha, member.ID, admins)
		err := tuples.Add(ctx, tuple)
		require.NoError(t, err)
	})
}

func TestTuplesRemove(t *testing.T) {
	ctx := context.Background()
	cleanup(conn)
	tuples := NewTuples(conn)
	objects := NewObjects(conn)
	roles := NewRoles(conn)

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "member",
		Verbs: []doorman.Verb{"foo"},
	}
	admins := doorman.Object("group:admins")

	t.Run("Fails if does not exist", func(t *testing.T) {
		err := tuples.Remove(ctx, doorman.NewTuple(alice, member.ID, admins))
		assert.Error(t, err)
	})

	// Make objects and roles
	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, objects.Add(ctx, admins))

	// Create tuple
	tuple := doorman.NewTuple(alice, member.ID, admins)
	require.NoError(t, tuples.Add(ctx, tuple))

	t.Run("Fails if does not exist", func(t *testing.T) {
		err := tuples.Remove(ctx, tuple)
		assert.NoError(t, err)
	})
}

func TestTuplesListConnected(t *testing.T) {
	ctx := context.Background()
	cleanup(conn)
	tuples := NewTuples(conn)
	objects := NewObjects(conn)
	roles := NewRoles(conn)

	// Init objects
	alice := doorman.Object("user:alice")
	users := doorman.Object("group:users")
	admins := doorman.Object("group:admins")
	superadmins := doorman.Object("group:superadmins")
	cia := doorman.Object("group:cia")

	// Init roles
	member := doorman.NewRole("member", []doorman.Verb{"foo"})

	// Create objects, roles
	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, users))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, superadmins))
	require.NoError(t, objects.Add(ctx, cia))
	require.NoError(t, roles.Add(ctx, member))

	// Create tuples
	require.NoError(t, tuples.Add(ctx, doorman.NewTuple(alice, member.ID, superadmins)))
	require.NoError(t, tuples.Add(ctx, doorman.NewTuple(cia, member.ID, superadmins)))
	require.NoError(t, tuples.Add(ctx, doorman.NewTuple(superadmins, member.ID, admins)))
	require.NoError(t, tuples.Add(ctx, doorman.NewTuple(admins, member.ID, users)))

	t.Run("ListConnected(alice, false)", func(t *testing.T) {
		paths, err := tuples.ListConnected(ctx, alice, false)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(paths))
		assert.Equal(t, doorman.Path{doorman.Connection{Role: member.ID, Object: superadmins}}, paths[0])
		assert.Equal(t, doorman.Path{
			doorman.Connection{Role: member.ID, Object: superadmins},
			doorman.Connection{Role: member.ID, Object: admins},
		}, paths[1])
		assert.Equal(t, doorman.Path{
			doorman.Connection{Role: member.ID, Object: superadmins},
			doorman.Connection{Role: member.ID, Object: admins},
			doorman.Connection{Role: member.ID, Object: users},
		}, paths[2])
	})

	t.Run("ListConnected(alice, true)", func(t *testing.T) {
		paths, err := tuples.ListConnected(ctx, alice, true)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(paths))
	})

	t.Run("ListConnected(alice, true)", func(t *testing.T) {
		paths, err := tuples.ListConnected(ctx, admins, true)
		assert.NoError(t, err)

		assert.Equal(t, 3, len(paths))
		assert.Equal(t, doorman.Path{doorman.Connection{Role: member.ID, Object: superadmins}}, paths[0])
		assert.Equal(t, doorman.Path{
			doorman.Connection{Role: member.ID, Object: superadmins},
			doorman.Connection{Role: member.ID, Object: cia},
		}, paths[1])
		assert.Equal(t, doorman.Path{
			doorman.Connection{Role: member.ID, Object: superadmins},
			doorman.Connection{Role: member.ID, Object: alice},
		}, paths[2])
	})
}
