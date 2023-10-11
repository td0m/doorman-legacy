package entities

import (
	"context"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/poc-doorman/entities/db"
)

func TestCreateType(t *testing.T) {
	ctx := context.Background()

	t.Run("Success on valid id", func(t *testing.T) {
		in := CreateTypeRequest{
			ID:    xid.New().String(),
			Attrs: map[string]any{"foo": "ba"},
		}

		res, err := CreateType(ctx, in)
		assert.NoError(t, err)
		assert.Equal(t, in.ID, res.ID)
		assert.Equal(t, in.Attrs, res.Attrs)
	})

	t.Run("Fails on invalid id", func(t *testing.T) {
		in := CreateTypeRequest{
			ID: xid.New().String() + "/blabla",
		}

		_, err := CreateType(ctx, in)
		assert.Error(t, err)
	})
}

func TestUpdateType(t *testing.T) {
	ctx := context.Background()
	dbtype := &db.Type{ID: xid.New().String()}
	require.NoError(t, dbtype.Create(ctx))

	in := UpdateTypeRequest{
		ID:    dbtype.ID,
		Attrs: map[string]any{"foo": "bar"},
	}

	res, err := UpdateType(ctx, in)
	assert.NoError(t, err)

	assert.Equal(t, dbtype.ID, res.ID)
	assert.Equal(t, in.Attrs, res.Attrs)

	updateddbtype, err := db.RetrieveType(ctx, dbtype.ID)
	require.NoError(t, err)
	assert.Equal(t, in.Attrs, updateddbtype.Attrs)
	assert.True(t, updateddbtype.UpdatedAt.After(dbtype.UpdatedAt))
}

