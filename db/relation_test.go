package db

import (
	"context"
	"fmt"
	"sync"
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

	ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
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

	ab := &Relation{ID: "ab_" + seed, From: a.ID, To: a.ID}
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

	ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, ab.ID)

	bc := &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
	fmt.Println("now")
	err = bc.Create(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, bc.ID)

	t.Run("Caches 'ab'", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select id from cache where "from"=$1 and "to"=$2`, a.ID, b.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
		require.Equal(t, "cache: "+ab.ID, rows[0]["id"])
	})

	t.Run("Caches 'ac'", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select id from cache where "from"=$1 and "to"=$2`, a.ID, c.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
		require.Equal(t, "cache: "+ab.ID+" "+bc.ID, rows[0]["id"])
	})

	t.Run("Does not cache 'ad' until after 'cd' added", func(t *testing.T) {
		rows := []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, d.ID)
		require.NoError(t, err)
		require.Equal(t, 0, len(rows))

		cd := &Relation{ID: "cd_" + seed, From: c.ID, To: d.ID}
		err = cd.Create(ctx)
		assert.NoError(t, err)

		rows = []map[string]any{}
		err = pgxscan.Select(ctx, pg, &rows, `select 1 from cache where "from"=$1 and "to"=$2`, a.ID, d.ID)
		require.NoError(t, err)
		require.Equal(t, 1, len(rows))
	})
}

func TestListRelationsRecHandlesCycles(t *testing.T) {
	t.Run("Left to Right", func(t *testing.T) {
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

		ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
		err = ab.Create(ctx)
		require.NoError(t, err)

		bc := &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
		err = bc.Create(ctx)
		require.NoError(t, err)

		// ca cannot be, as we already have ac (although implicitly)
		ca := &Relation{ID: "ca_" + seed, From: c.ID, To: a.ID}
		err = ca.Create(ctx)
		assert.ErrorIs(t, err, ErrCycle)
	})

	t.Run("Right To Left", func(t *testing.T) {
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

		ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
		err = ab.Create(ctx)
		require.NoError(t, err)

		// ca cannot be, as we already have ac (although implicitly)
		ca := &Relation{ID: "ca_" + seed, From: c.ID, To: a.ID}
		err = ca.Create(ctx)
		require.NoError(t, err)

		bc := &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
		err = bc.Create(ctx)
		assert.ErrorIs(t, err, ErrCycle)

	})
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

	ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	require.NoError(t, err)

	a2b := &Relation{ID: "a2b_" + seed, From: a2.ID, To: b.ID}
	err = a2b.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
	err = bc.Create(ctx)
	require.NoError(t, err)

	cd := &Relation{ID: "cd_" + seed, From: c.ID, To: d.ID}
	err = cd.Create(ctx)
	require.NoError(t, err)

	de := &Relation{ID: "de_" + seed, From: d.ID, To: e.ID}
	err = de.Create(ctx)
	require.NoError(t, err)

	t.Run("To 'd'", func(t *testing.T) {
		tx, err := pg.Begin(ctx)
		require.NoError(t, err)

		defer func() {
			_ = tx.Commit(ctx)
		}()

		rels, err := listRecRelationsTo(ctx, tx, d.ID)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(rels))

		assert.Equal(t, cd.ID, rels[0].ID)
		assert.Equal(t, []string{cd.ID}, rels[0].Via)

		assert.Equal(t, bc.ID, rels[1].ID)
		assert.Equal(t, []string{cd.ID, bc.ID}, rels[1].Via)

		assert.Equal(t, ab.ID, rels[2].ID)
		assert.Equal(t, []string{cd.ID, bc.ID, ab.ID}, rels[2].Via)

		assert.Equal(t, a2b.ID, rels[3].ID)
		assert.Equal(t, []string{cd.ID, bc.ID, a2b.ID}, rels[3].Via)
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
		assert.Equal(t, []string{cd.ID}, rels[0].Via)

		assert.Equal(t, de.ID, rels[1].ID)
		assert.Equal(t, []string{cd.ID, de.ID}, rels[1].Via)
	})
}

func TestRelationsDelete(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	// Entities

	a := &Entity{ID: "user:a_" + seed}
	err := a.Create(ctx)
	require.NoError(t, err)

	b := &Entity{ID: "user:b_" + seed}
	err = b.Create(ctx)
	require.NoError(t, err)

	c := &Entity{ID: "user:c_" + seed}
	err = c.Create(ctx)
	require.NoError(t, err)

	// Relations

	ab := &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
	err = ab.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
	err = bc.Create(ctx)
	require.NoError(t, err)

	// Caches

	abcaches, err := ListRelationsOrCache(ctx, "cache", RelationFilter{
		From: &a.ID,
		To:   &b.ID,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(abcaches))

	accaches, err := ListRelationsOrCache(ctx, "cache", RelationFilter{
		From: &a.ID,
		To:   &c.ID,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(accaches))

	// Delete

	err = ab.Delete(ctx)
	require.NoError(t, err)

	// Caches after deletion

	abcaches, err = ListRelationsOrCache(ctx, "cache", RelationFilter{
		From: &a.ID,
		To:   &b.ID,
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(abcaches))

	accaches, err = ListRelationsOrCache(ctx, "cache", RelationFilter{
		From: &a.ID,
		To:   &c.ID,
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(accaches))
}

func TestCreateRelationsInParallel(t *testing.T) {
	ctx := context.Background()
	seed := xid.New().String()

	// Entities

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

	// Relations
	var ab *Relation
	var a2b *Relation
	var bc *Relation
	var cd *Relation
	var de *Relation

	funcs := []func(){
		func() {
			ab = &Relation{ID: "ab_" + seed, From: a.ID, To: b.ID}
			err:= ab.Create(ctx)
			require.NoError(t, err)
		},

		func() {
			a2b = &Relation{ID: "a2b_" + seed, From: a2.ID, To: b.ID}
			err:= a2b.Create(ctx)
			require.NoError(t, err)
		},

		func() {
			bc = &Relation{ID: "bc_" + seed, From: b.ID, To: c.ID}
			err:= bc.Create(ctx)
			require.NoError(t, err)
		},

		func() {
			cd = &Relation{ID: "cd_" + seed, From: c.ID, To: d.ID}
			err:= cd.Create(ctx)
			require.NoError(t, err)
		},

		func() {
			de = &Relation{ID: "de_" + seed, From: d.ID, To: e.ID}
			err:= de.Create(ctx)
			require.NoError(t, err)
		},
	}

	var wg sync.WaitGroup
	wg.Add(len(funcs))
	for i := range funcs {
		i := i
		go func() {
			defer wg.Done()
			funcs[i]()
		}()
	}

	wg.Wait()

	t.Run("To 'd'", func(t *testing.T) {
		tx, err := pg.Begin(ctx)
		require.NoError(t, err)

		defer func() {
			_ = tx.Commit(ctx)
		}()

		rels, err := listRecRelationsTo(ctx, tx, d.ID)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(rels))

		assert.Equal(t, cd.ID, rels[0].ID)
		assert.Equal(t, []string{cd.ID}, rels[0].Via)

		assert.Equal(t, bc.ID, rels[1].ID)
		assert.Equal(t, []string{cd.ID, bc.ID}, rels[1].Via)

		assert.Equal(t, ab.ID, rels[2].ID)
		assert.Equal(t, []string{cd.ID, bc.ID, ab.ID}, rels[2].Via)

		assert.Equal(t, a2b.ID, rels[3].ID)
		assert.Equal(t, []string{cd.ID, bc.ID, a2b.ID}, rels[3].Via)
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
		assert.Equal(t, []string{cd.ID}, rels[0].Via)

		assert.Equal(t, de.ID, rels[1].ID)
		assert.Equal(t, []string{cd.ID, de.ID}, rels[1].Via)
	})
}
