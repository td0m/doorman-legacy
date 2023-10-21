package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanup(t *testing.T) {
	ctx := context.Background()
	_, err := pg.Exec(ctx, `delete from dependencies;`)
	if err != nil {
		t.Fail()
	}
	_, err = pg.Exec(ctx, `delete from relations;`)
	if err != nil {
		t.Fail()
	}
}

func newRelation(from, name, to string) *Relation {
	return &Relation{
		From: from,
		Name: name,
		To:   to,
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pool, err := pgxpool.New(ctx, "user=doorman database=doorman")
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	pg = pool

	m.Run()
}

func TestRelationCreateSuccess(t *testing.T) {
	fmt.Println(context.Background())
	ctx := context.Background()

	r := &Relation{
		From: "repo:" + xid.New().String(),
		To:   "user:" + xid.New().String(),
	}
	err := r.Create(ctx)
	assert.NoError(t, err)

	var saved Relation
	err = pg.QueryRow(ctx, `select "from", "to" from relations where id=$1`, r.ID).Scan(&saved.From, &saved.To)
	assert.NoError(t, err)

	assert.Equal(t, r.From, saved.From)
	assert.Equal(t, r.To, saved.To)
}

func TestRelationCreateFailureOnCycle(t *testing.T) {
}

func TestRelationCreateTypeValidation(t *testing.T) {
	// always needs type
	// user -> *
	// user -> collection
	// anything not a user -> collection
}

func TestCollectionCopiesUser(t *testing.T) {
	cleanup(t)

	ctx := context.Background()

	// TODO: check all
	// 1. user added to collection
	// 2. new parent collection added
	// 3. connected resource to collection

	// setup relations

	fooOwnerAdmins := newRelation("resource:foo", "owner", "collection:admins")
	err := fooOwnerAdmins.Create(ctx)
	require.NoError(t, err)

	adminsMemberAlice := newRelation("collection:admins", "member", "user:alice")
	err = adminsMemberAlice.Create(ctx)
	require.NoError(t, err)

	adminsChildSuperadmins := newRelation("collection:admins", "child", "collection:superadmins")
	err = adminsChildSuperadmins.Create(ctx)
	require.NoError(t, err)

	superadminsMemberBob := newRelation("collection:superadmins", "member", "user:bob")
	err = superadminsMemberBob.Create(ctx)
	require.NoError(t, err)

	var relations []Relation
	err = pgxscan.Select(ctx, pg, &relations, `select id, "to", name from relations where "from"=$1`, "resource:foo")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(relations))

	assert.Equal(t, "owner", relations[0].Name)
	assert.Equal(t, "collection:admins", relations[0].To)

	{
		assert.Equal(t, "owner", relations[1].Name)
		assert.Equal(t, "user:alice", relations[1].To)
		var deps []Dependency
		err = pgxscan.Select(ctx, pg, &deps, `select depends_on from dependencies where relation_id=$1`, relations[1].ID)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(deps))
		assert.Equal(t, fooOwnerAdmins.ID, deps[0].DependsOn)
	}

	{
		assert.Equal(t, "owner", relations[2].Name)
		assert.Equal(t, "user:bob", relations[2].To)
		var deps []Dependency
		err = pgxscan.Select(ctx, pg, &deps, `select depends_on from dependencies where relation_id=$1`, relations[2].ID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(deps))
		assert.Equal(t, adminsChildSuperadmins.ID, deps[0].DependsOn)
		assert.Equal(t, fooOwnerAdmins.ID, deps[1].DependsOn)
	}
}

// TODO: ensure gets deleted if dependency also deleted

func TestListRelationsRec(t *testing.T) {
	ctx := context.Background()

	// Create Relations

	ab := &Relation{ID: "ab", From: "collection:a", To: "collection:b"}
	err := ab.Create(ctx)
	require.NoError(t, err)

	a2b := &Relation{ID: "a2b", From: "collection:a2", To: "collection:b"}
	err = a2b.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{ID: "bc", From: "collection:b", To: "collection:c"}
	err = bc.Create(ctx)
	require.NoError(t, err)

	cd := &Relation{ID: "cd", From: "collection:c", To: "collection:d"}
	err = cd.Create(ctx)
	require.NoError(t, err)

	de := &Relation{ID: "de", From: "collection:d", To: "collection:e"}
	err = de.Create(ctx)
	require.NoError(t, err)

	t.Run("To 'd'", func(t *testing.T) {
		tx, err := pg.Begin(ctx)
		require.NoError(t, err)

		defer func() {
			_ = tx.Commit(ctx)
		}()

		rels, err := listRecRelationsTo(ctx, tx, "collection:d")
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

		rels, err := listRecRelationsFrom(ctx, tx, "collection:c")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(rels))
		assert.Equal(t, "cd", rels[0].ID)
		assert.Equal(t, []string{"cd"}, rels[0].Via)

		assert.Equal(t, "de", rels[1].ID)
		assert.Equal(t, []string{"cd", "de"}, rels[1].Via)
	})
}
