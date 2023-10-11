package entities

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/u"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)
	defer pgDoorman.Close()

	db.Conn = pgDoorman

	m.Run()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		in      CreateRequest
		success bool
	}{
		// {CreateRequest{Type: ""}, false},
		// {CreateRequest{Type: "no spaces"}, false},
		// {CreateRequest{Type: "bad_characters"}, false},
		// {CreateRequest{Type: "0numbers"}, false},
		// {CreateRequest{Type: "n0numbers"}, true},
		// {CreateRequest{Type: "LOWERCASE"}, false},

		{CreateRequest{Type: "user"}, true},
		{CreateRequest{Type: "user", Attrs: map[string]any{}}, true},
		{CreateRequest{Type: "user", Attrs: map[string]any{"foo": "bar"}}, true},

		{CreateRequest{Type: "user", ID: "path/to/resources/128"}, false},
		{CreateRequest{Type: "user", ID: "no spaces"}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Entity: %+v", tt.in), func(t *testing.T) {
			_ = (&db.Entity{ID: tt.in.ID, Type: tt.in.Type}).Delete(ctx)
			res, err := Create(ctx, tt.in)

			if tt.success {
				require.NoError(t, err)
				require.Equal(t, tt.in.Type, res.Type)
				if len(tt.in.Attrs) > 0 {
					require.Equal(t, tt.in.Attrs, res.Attrs)
				}
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	dbentity := &db.Entity{
		Type: "user",
	}
	err := dbentity.Create(ctx)
	require.NoError(t, err)

	in := UpdateRequest{
		ID:    dbentity.ID,
		Type:  dbentity.Type,
		Attrs: map[string]any{"bar": true},
	}
	res, err := Update(ctx, in)
	require.NoError(t, err)

	require.Equal(t, in.Attrs, res.Attrs)
	require.True(t, res.UpdatedAt.After(res.CreatedAt))

	dbentity, err = db.Get(ctx, in.ID, in.Type)
	require.NoError(t, err)

	require.Equal(t, in.Attrs, dbentity.Attrs)
	require.NotNil(t, dbentity.UpdatedAt)
	require.True(t, dbentity.UpdatedAt.After(dbentity.CreatedAt))
}

func TestList(t *testing.T) {
	ctx := context.Background()

	rnd := xid.New().String()

	dbtype := &db.Type{ID: rnd}
	require.NoError(t, dbtype.Create(ctx))

	ids := make([]string, 1_200)
	for i := range ids {
		ids[i] = xid.New().String()
		dbentity := &db.Entity{ID: ids[i], Type: rnd}
		err := dbentity.Create(ctx)
		require.NoError(t, err)
	}

	t.Run("Only returns 1000 first items without pagination token", func(t *testing.T) {
		res, err := List(ctx, ListRequest{
			Type: &rnd,
		})
		require.NoError(t, err)
		assert.Equal(t, 1_000, len(res.Data))
		assert.Equal(t, ids[999], res.Data[999].ID)
		lastToken := res.PaginationToken
		t.Run("Second request returns remaining 200 items", func(t *testing.T) {
			res, err := List(ctx, ListRequest{
				Type:            &rnd,
				PaginationToken: lastToken,
			})
			require.NoError(t, err)
			assert.Equal(t, 200, len(res.Data))
			assert.Equal(t, ids[1199], res.Data[199].ID)
		})
	})
}

