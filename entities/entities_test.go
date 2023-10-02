package entities

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
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
		in      Entity
		success bool
	}{
		{Entity{Type: ""}, false},
		{Entity{Type: "foo"}, true},
		{Entity{Type: "foo", Attrs: map[string]any{}}, true},
		{Entity{Type: "foo", Attrs: map[string]any{"foo": "bar"}}, true},
		{Entity{Type: "no spaces"}, false},
		{Entity{Type: "bad_characters"}, false},
		{Entity{Type: "n0numbers"}, false},
		{Entity{Type: "LOWERCASE"}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Entity: %+v", tt.in), func(t *testing.T) {
			res, err := Create(ctx, tt.in)

			if tt.success {
				require.NoError(t, err)
				require.Equal(t, tt.in.Type, res.Type)
				if tt.in.Attrs != nil {
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
		Type: "foo",
	}
	err := dbentity.Create(ctx)
	require.NoError(t, err)

	in := UpdateRequest{
		ID: dbentity.ID,
		Type: dbentity.Type,
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
