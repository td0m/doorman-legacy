package db

import (
	"context"
	"testing"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/poc-doorman/u"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pool, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)

	defer pool.Close()

	pg = pool

	m.Run()
}

func TestRelationCreateSuccess(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	b := &Entity{ID: "user:b_" + seed}
	err = b.Create(ctx)
	require.NoError(t, err)

	ab := &Relation{From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, ab.ID)

	row := map[string]any{}
	err = get(ctx, &row, `select "from", "to" from relations where id=$1`, ab.ID)
	require.NoError(t, err)
	require.Equal(t, ab.From, row["from"])
	require.Equal(t, ab.To, row["to"])
}

func TestRelationCreateFailureToConnectToSelf(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	ab := &Relation{From: a.ID, To: a.ID}
	err = ab.Create(ctx)
	assert.Error(t, err)
}

// TODO: validate based on types as well

func TestRelationCreateBuildsCache(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	b := &Entity{ID: "user:b_" + seed}
	err = b.Create(ctx)
	require.NoError(t, err)

	c := &Entity{ID: "user:c_" + seed}
	err = c.Create(ctx)
	require.NoError(t, err)

	d := &Entity{ID: "user:d_" + seed}
	err = d.Create(ctx)
	require.NoError(t, err)

	ab := &Relation{From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, ab.ID)

	bc := &Relation{From: b.ID, To: c.ID}
	err = bc.Create(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, bc.ID)

	t.Run("Caches 'ab'", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, b.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
	})

	t.Run("Caches 'ac'", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, c.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
	})

	t.Run("Does not cache 'ad' until after 'cd' added", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, d.ID)
		require.NoError(t, err)
		require.Equal(t, 0, len(rows))

		cd := &Relation{From: c.ID, To: d.ID}
		err = cd.Create(ctx)
		assert.NoError(t, err)

		rows = []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, d.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
	})
}

func TestListRelationsRecHandlesCycles(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	b := &Entity{ID: "user:b_" + seed}
	err = b.Create(ctx)
	require.NoError(t, err)

	c := &Entity{ID: "user:c_" + seed}
	err = c.Create(ctx)
	require.NoError(t, err)

	ab := &Relation{From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{From: b.ID, To: c.ID}
	err = bc.Create(ctx)
	require.NoError(t, err)

	// ca cannot be!
}

func TestListRelationsRec(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	// Create Entities

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	a2 := &Entity{ID: "user:a2_" + seed}
	err = a2.Create(ctx)
	require.NoError(t, err)

	b := &Entity{ID: "user:b_" + seed}
	err = b.Create(ctx)
	require.NoError(t, err)

	c := &Entity{ID: "user:c_" + seed}
	err = c.Create(ctx)
	require.NoError(t, err)

	d := &Entity{ID: "user:d_" + seed}
	err = d.Create(ctx)
	require.NoError(t, err)

	e := &Entity{ID: "user:e_" + seed}
	err = e.Create(ctx)
	require.NoError(t, err)

	// Create Relations

	ab := &Relation{From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	require.NoError(t, err)

	a2b := &Relation{From: a2.ID, To: b.ID}
	err = a2b.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{From: b.ID, To: c.ID}
	err = bc.Create(ctx)
	require.NoError(t, err)

	cd := &Relation{From: c.ID, To: d.ID}
	err = cd.Create(ctx)
	require.NoError(t, err)

	de := &Relation{From: d.ID, To: e.ID}
	err = de.Create(ctx)
	require.NoError(t, err)

	t.Run("To 'c'", func(t *testing.T) {
		tx, err := pg.Begin(ctx)
		require.NoError(t, err)

		defer func() {
			_ = tx.Commit(ctx)
		}()

		rels, err := listRecRelationsTo(ctx, tx, c.ID)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(rels))

		assert.Equal(t, bc.ID, rels[0].ID)
		assert.Equal(t, []string{}, rels[0].Via)

		assert.Equal(t, ab.ID, rels[1].ID)
		assert.Equal(t, []string{bc.ID}, rels[1].Via)

		assert.Equal(t, a2b.ID, rels[2].ID)
		assert.Equal(t, []string{bc.ID}, rels[2].Via)
	})

	t.Run("From 'c'", func(t *testing.T) {
		tx, err := pg.Begin(ctx)
		require.NoError(t, err)

		defer func() {
			_ = tx.Commit(ctx)
		}()

		rels, err := listRecRelationsFrom(ctx, tx, c.ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(rels))
		assert.Equal(t, cd.ID, rels[0].ID)
		assert.Equal(t, []string{}, rels[0].Via)

		assert.Equal(t, de.ID, rels[1].ID)
		assert.Equal(t, []string{cd.ID}, rels[1].Via)
	})
}
