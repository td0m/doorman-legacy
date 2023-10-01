package relations

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	entitiesdb "github.com/td0m/poc-doorman/entities/db"
	relationsdb "github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)
	defer pgDoorman.Close()

	entitiesdb.Conn = pgDoorman
	relationsdb.Conn = pgDoorman

	m.Run()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	t.Run("Success on valid entities", func(t *testing.T) {
		user := &entitiesdb.Entity{
			Type: "user",
			ID:   xid.New().String(),
		}
		require.NoError(t, user.Create(ctx))

		resource := &entitiesdb.Entity{
			Type: "user",
			ID:   xid.New().String(),
		}
		require.NoError(t, resource.Create(ctx))

		req := CreateRequest{
			ID:    xid.New().String(),
			From:  Entity{ID: user.ID, Type: user.Type},
			To:    Entity{ID: resource.ID, Type: resource.Type},
			Attrs: map[string]any{},
		}
		relation, err := Create(ctx, req)
		require.NoError(t, err)
		require.Equal(t, req.ID, relation.ID)
		require.Equal(t, req.From.ID, relation.From.ID)
		require.Equal(t, req.From.Type, relation.From.Type)
		require.Equal(t, req.To.ID, relation.To.ID)
		require.Equal(t, req.To.Type, relation.To.Type)
		require.Equal(t, req.Attrs, relation.Attrs)
	})

	t.Run("Failure on missing \"from\" entity", func(t *testing.T) {
		user := &entitiesdb.Entity{
			Type: "user",
			ID:   xid.New().String(),
		}
		require.NoError(t, user.Create(ctx))

		req := CreateRequest{
			ID:   xid.New().String(),
			From: Entity{ID: user.ID, Type: user.Type},
		}
		_, err := Create(ctx, req)
		require.Error(t, err)
	})

	t.Run("Fails on cyclic relation", func(t *testing.T) {
		// TODO
	})
}
