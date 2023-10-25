package server

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
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

func TestCheckDirect(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations()
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "member",
		Verbs: []doorman.Verb{"foo"},
	}
	admins := doorman.Object("group:admins")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, roles.Add(ctx, member))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "foo",
			Object:  string(admins),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "foo",
			Object:  string(admins),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}
