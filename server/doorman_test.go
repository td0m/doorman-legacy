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
		delete from changes;
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	admins := doorman.Object("group:admins")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.roles.Add(ctx, member))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "inherits",
			Object:  string(admins),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "inherits",
			Object:  string(admins),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckViaGroup(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, member))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, superadmins))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, member))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(superadmins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	duperadmins := doorman.Object("group:duperadmins")
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, duperadmins))
	require.NoError(t, s.objects.Add(ctx, superadmins))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, member))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(duperadmins),
			Role:    member.ID,
			Object:  string(superadmins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(duperadmins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)

	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "group:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	duperadmins := doorman.Object("group:duperadmins")
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, duperadmins))
	require.NoError(t, s.objects.Add(ctx, superadmins))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, member))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(alice),
				Role:    member.ID,
				Object:  string(duperadmins),
			})
			return err
		})

		g.Go(func() error {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(duperadmins),
				Role:    member.ID,
				Object:  string(superadmins),
			})
			return err
		})

		g.Go(func() error {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(superadmins),
				Role:    member.ID,
				Object:  string(admins),
			})
			return err
		})

		g.Go(func() error {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(admins),
				Role:    owner.ID,
				Object:  string(banana),
			})
			return err
		})

		require.NoError(t, g.Wait())
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	member := doorman.Role{
		ID:    "groop:member",
		Verbs: []doorman.Verb{"inherits"},
	}
	admins := doorman.Object("groop:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, member))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    member.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"inherits"}}
	groopMember := doorman.Role{ID: "groop:member", Verbs: []doorman.Verb{"inherits"}}
	superadmins := doorman.Object("group:superadmins")
	admins := doorman.Object("groop:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, superadmins))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, groupMember))
	require.NoError(t, s.roles.Add(ctx, groopMember))
	require.NoError(t, s.roles.Add(ctx, owner))

	t.Run("Failure: Check before granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})

	t.Run("Grant", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(admins),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(superadmins),
			Role:    groopMember.ID,
			Object:  string(admins),
		})
		require.NoError(t, err)

		_, err = s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    groupMember.ID,
			Object:  string(superadmins),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after granting", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
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

	s := NewDoorman(conn)
	ctx := context.Background()

	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"inherits"}}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, groupMember))
	require.NoError(t, s.roles.Add(ctx, owner))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(admins),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Fails directly", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(banana),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.ErrorIs(t, err, db.ErrCycle)
	})

	t.Run("Fails indirectly", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(banana),
			Role:    groupMember.ID,
			Object:  string(admins),
		})
		require.ErrorIs(t, err, db.ErrCycle)
	})
}

func TestConnectingToSelfIndirectlyInParallelFails(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"inherits"}}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, groupMember))
	require.NoError(t, s.roles.Add(ctx, owner))

	// run many times as failing cound happen at random
	for i := 0; i < 100; i++ {
		_, _ = conn.Exec(ctx, `delete from tuples`)

		var g errgroup.Group
		g.Go(func() error {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(banana),
				Role:    groupMember.ID,
				Object:  string(admins),
			})
			return err
		})
		g.Go(func() error {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: string(admins),
				Role:    owner.ID,
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

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, owner))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Success: Check before revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})

	t.Run("Successful revoke of existing role", func(t *testing.T) {
		_, err = s.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(alice),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Faiure revoking same role twice (non-idempotent)", func(t *testing.T) {
		_, err = s.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(alice),
			Role:    owner.ID,
			Object:  string(banana),
		})
		require.ErrorIs(t, err, doorman.ErrTupleNotFound)
	})

	t.Run("Failure: Check after revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})
}

func TestCheckRevokeOneOfTwo(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	eater := doorman.Role{
		ID:    "item:eater",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, owner))
	require.NoError(t, s.roles.Add(ctx, eater))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	_, err = s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    eater.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Success: Check before revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})

	t.Run("Successful revoke of one of two s.roles", func(t *testing.T) {
		_, err = s.Revoke(ctx, &pb.RevokeRequest{
			Subject: string(alice),
			Role:    eater.ID,
			Object:  string(banana),
		})
		require.NoError(t, err)
	})

	t.Run("Success: Check after revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})
}

func TestCheckUpdateRole(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, owner))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Check before revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Success)

		res, err = s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "drink",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Success)
	})

	_, err = s.UpsertRole(ctx, &pb.UpsertRoleRequest{
		Id:    owner.ID,
		Verbs: []string{"drink"},
	})
	require.NoError(t, err)

	t.Run("Check after revoking", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, false, res.Success)

		res, err = s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "drink",
			Object:  string(banana),
		})
		assert.NoError(t, err)
		assert.Equal(t, true, res.Success)
	})
}

func TestCheckUpsertRoleCreate(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, banana))

	t.Run("Failed to grant before role is created", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "owner",
			Object:  string(banana),
		})
		require.Error(t, err)
	})

	_, err := s.UpsertRole(ctx, &pb.UpsertRoleRequest{
		Id:    "item:owner",
		Verbs: []string{"drink"},
	})
	require.NoError(t, err)

	t.Run("Success grant after role is created", func(t *testing.T) {
		_, err := s.Grant(ctx, &pb.GrantRequest{
			Subject: string(alice),
			Role:    "item:owner",
			Object:  string(banana),
		})
		require.NoError(t, err)
	})
}

func TestRemoveRole(t *testing.T) {
	cleanup(conn)

	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.roles.Add(ctx, owner))
	require.NoError(t, s.objects.Add(ctx, banana))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	t.Run("Can access role before removing", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, true, res.Success)
	})

	t.Run("Success removing role", func(t *testing.T) {
		_, err := s.RemoveRole(ctx, &pb.RemoveRoleRequest{
			Id: "item:owner",
		})
		require.NoError(t, err)
	})

	t.Run("Fails to remove role again", func(t *testing.T) {
		_, err := s.RemoveRole(ctx, &pb.RemoveRoleRequest{
			Id: "item:owner",
		})
		require.Error(t, err)
	})

	t.Run("Removed existing tuples with the role", func(t *testing.T) {
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.Equal(t, false, res.Success)
	})
}

// func TestListChanges(t *testing.T) {
// 	cleanup(conn)
// 	s := NewDoorman(conn)
// 	ctx := context.Background()
//
// 	alice := doorman.Object("user:alice")
// 	owner := doorman.Role{
// 		ID:    "item:owner",
// 		Verbs: []doorman.Verb{"eat"},
// 	}
// 	banana := doorman.Object("item:banana")
//
// 	require.NoError(t, s.objects.Add(ctx, alice))
// 	require.NoError(t, s.roles.Add(ctx, owner))
// 	require.NoError(t, s.objects.Add(ctx, banana))
//
// 	{
// 		_, err := s.Grant(ctx, &pb.GrantRequest{
// 			Subject: string(alice),
// 			Role:    string("owner"),
// 			Object:  string(banana),
// 		})
// 		require.NoError(t, err)
// 	}
//
// 	res, err := s.Changes(ctx, &pb.ChangesRequest{})
// 	require.NoError(t, err)
// 	require.Equal(t, 2, len(res.Items))
// 	require.Equal(t, "TUPLE_CREATED", res.Items[0].Type)
// 	require.Equal(t, "RELATION_CREATED", res.Items[1].Type)
//
// 	{
// 		_, err := s.Revoke(ctx, &pb.RevokeRequest{
// 			Subject: string(alice),
// 			Role:    string("owner"),
// 			Object:  string(banana),
// 		})
// 		require.NoError(t, err)
// 	}
//
// 	res, err = s.Changes(ctx, &pb.ChangesRequest{PaginationToken: res.PaginationToken})
// 	require.NoError(t, err)
// 	require.Equal(t, 2, len(res.Items))
// 	require.Equal(t, "TUPLE_REMOVED", res.Items[0].Type)
// 	require.Equal(t, "RELATION_REMOVED", res.Items[1].Type)
// }

func TestRemoveOneOfTwoRolesWithSameVerb(t *testing.T) {
	cleanup(conn)
	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	reader := doorman.Role{
		ID:    "item:reader",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.roles.Add(ctx, owner))
	require.NoError(t, s.roles.Add(ctx, reader))
	require.NoError(t, s.objects.Add(ctx, banana))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)
	_, err = s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    reader.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.True(t, res.Success)
	}

	_, err = s.Revoke(ctx, &pb.RevokeRequest{
		Subject: string(alice),
		Role:    reader.ID,
		Object:  string(banana),
	})
	assert.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.True(t, res.Success)
	}

	_, err = s.Revoke(ctx, &pb.RevokeRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	assert.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.False(t, res.Success)
	}
}

func TestRemoveOneOfTwoRolesWithSameVerb2(t *testing.T) {
	cleanup(conn)
	s := NewDoorman(conn)
	ctx := context.Background()

	alice := doorman.Object("user:alice")
	groupMember := doorman.Role{ID: "group:member", Verbs: []doorman.Verb{"inherits"}}
	admins := doorman.Object("group:admins")
	owner := doorman.Role{
		ID:    "item:owner",
		Verbs: []doorman.Verb{"eat"},
	}
	banana := doorman.Object("item:banana")

	require.NoError(t, s.objects.Add(ctx, alice))
	require.NoError(t, s.objects.Add(ctx, admins))
	require.NoError(t, s.objects.Add(ctx, banana))
	require.NoError(t, s.roles.Add(ctx, groupMember))
	require.NoError(t, s.roles.Add(ctx, owner))

	_, err := s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    groupMember.ID,
		Object:  string(admins),
	})
	require.NoError(t, err)

	_, err = s.Grant(ctx, &pb.GrantRequest{
		Subject: string(admins),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	_, err = s.Grant(ctx, &pb.GrantRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	require.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.True(t, res.Success)
	}

	_, err = s.Revoke(ctx, &pb.RevokeRequest{
		Subject: string(admins),
		Role:    owner.ID,
		Object:  string(banana),
	})
	assert.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.True(t, res.Success)
	}

	_, err = s.Revoke(ctx, &pb.RevokeRequest{
		Subject: string(alice),
		Role:    owner.ID,
		Object:  string(banana),
	})
	assert.NoError(t, err)

	{
		res, err := s.Check(ctx, &pb.CheckRequest{
			Subject: string(alice),
			Verb:    "eat",
			Object:  string(banana),
		})
		require.NoError(t, err)
		require.False(t, res.Success)
	}
}
