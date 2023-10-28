package server

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/doorman"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"golang.org/x/sync/errgroup"
)

var conn *pgxpool.Pool

func cleanup(conn *pgxpool.Pool) {
	ctx := context.Background()
	_, err := conn.Exec(ctx, `
		delete from tuples;
		delete from roles;
		delete from objects;
		delete from relations;
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

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
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
			Role:    "member",
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

func TestCheckViaGroup(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"foo"},
	}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "member",
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    "owner",
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckViaTwoGroups(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"foo"},
	}
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, superadmins))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "member",
			Object:  string(superadmins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    "owner",
			Object:  string(banana),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    "member",
			Object:  string(admins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckViaThreeGroups(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"foo"},
	}
	duperadmins := doorman.Object("group:duperadmins")
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, duperadmins))
	require.NoError(t, objects.Add(ctx, superadmins))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    "member",
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(duperadmins),
			Role:    "member",
			Object:  string(superadmins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "member",
			Object:  string(duperadmins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    "owner",
			Object:  string(banana),
		})
		require.NoError(t, err)

	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckViaThreeGroupsGrantedInParallel(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"foo"},
	}
	duperadmins := doorman.Object("group:duperadmins")
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, duperadmins))
	require.NoError(t, objects.Add(ctx, superadmins))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		var g errgroup.Group

		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(alice),
				Role:    "member",
				Object:  string(duperadmins),
			})
			return err
		})

		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(duperadmins),
				Role:    "member",
				Object:  string(superadmins),
			})
			return err
		})

		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(superadmins),
				Role:    "member",
				Object:  string(admins),
			})
			return err
		})

		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(admins),
				Role:    "owner",
				Object:  string(banana),
			})
			return err
		})

		require.NoError(t, g.Wait())
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckViaGroop(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "groop:member",
		Verbs: []doorman.Verb{"foo"},
	}
	admins := doorman.Object("groop:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, member))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "member",
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    "owner",
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})
}

func TestCheckViaGroupAndGroop(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"foo"}}
	groopMember := doorman.Role{ID: "groop:member", Verbs: []doorman.Verb{"foo"}}
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("groop:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, superadmins))
	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, groupMember))
	require.NoError(t, roles.Add(ctx, groopMember))
	require.NoError(t, roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    "owner",
			Object:  string(banana),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    "member",
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = server.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "member",
			Object:  string(superadmins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

}

func TestConnectingToSelfFails(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"foo"}}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, groupMember))
	require.NoError(t, roles.Add(ctx, owner))

	_, err := server.Grant(ctx, &pb.GrantRequest{
		Subject: string(admins),
		Role:    "owner",
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Fails directly", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(banana),
			Role:    "owner",
			Object:  string(banana),
		})
		require.ErrorIs(t, err, db.ErrCycle)
	})

	t.Run("Fails indirectly", func(t *testing.T) {
		_, err := server.Grant(ctx, &pb.GrantRequest{
			Subject: string(banana),
			Role:    "member",
			Object:  string(admins),
		})
		require.ErrorIs(t, err, db.ErrCycle)
	})
}

func TestConnectingToSelfIndirectlyInParallelFails(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"foo"}}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, admins))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, groupMember))
	require.NoError(t, roles.Add(ctx, owner))

	// run many times as failing cound happen at random
	for i := 0; i < 100; i++ {
		_, _ = conn.Exec(ctx, `delete from tuples`)

		var g errgroup.Group
		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(banana),
				Role:    "member",
				Object:  string(admins),
			})
			return err
		})
		g.Go(func() error {
			_, err := server.Grant(ctx, &pb.GrantRequest{
				Subject: string(admins),
				Role:    "owner",
				Object:  string(banana),
			})
			return err
		})

		err := g.Wait()
		require.ErrorIs(t, err, db.ErrCycle, "run %d", i)
	}
}

func TestCheckRevoke(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, owner))

	_, err := server.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    "owner",
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Success: Check before revoking", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})

	_, err = server.Revoke(ctx, &pb.RevokeRequest{
		Subject: string(alice),
		Role:    "owner",
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Failure: Check after revoking", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})
}

func TestCheckUpdateRole(t *testing.T) {
	cleanup(conn)

	relations := db.NewRelations(conn)
	objects := db.NewObjects(conn)
	roles := db.NewRoles(conn)
	tuples := db.NewTuples(conn)

	server := NewDoorman(relations, roles, tuples)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, objects.Add(ctx, alice))
	require.NoError(t, objects.Add(ctx, banana))
	require.NoError(t, roles.Add(ctx, owner))

	_, err := server.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    "owner",
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Check before revoking", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Success)

		res, err = server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "drink",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Success)
	})

	_, err = server.UpdateRole(ctx, &pb.UpdateRoleRequest{
		Id:    owner.ID,
		Verbs: []string{"drink"},
	})
	require.NoError(t, err)

	t.Run("Check after revoking", func(t *testing.T) {
		res, err := server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Success)

		res, err = server.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "drink",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Success)
	})
}
