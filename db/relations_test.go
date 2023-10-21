package db

import (
	"context"
	"testing"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func cleanup(t *testing.T) {
	ctx := context.Background()
	_, err := pg.Exec(ctx, `delete from relations;`)
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
	cleanup(t)
	ctx := context.Background()

	r := &Relation{
		From: "repo:foo",
		To:   "user:bar",
	}
	err := r.Create(ctx)
	assert.NoError(t, err)

	var saved Relation
	err = pg.QueryRow(ctx, `select "from", "to" from relations where "from"=$1 and "to"=$2`, r.From, r.To).Scan(&saved.From, &saved.To)
	assert.NoError(t, err)
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
	err = pgxscan.Select(ctx, pg, &relations, `select "to", name, via from relations where "from"=$1`, "resource:foo")
	assert.NoError(t, err)

	assert.Equal(t, 3, len(relations))

	assert.Equal(t, "owner", relations[0].Name)
	assert.Equal(t, "collection:admins", relations[0].To)

	assert.Equal(t, "owner", relations[1].Name)
	assert.Equal(t, "user:alice", relations[1].To)
	assert.Equal(t, []string{"resource:foo", "owner", "collection:admins", "member", "user:alice"}, relations[1].Via)

	assert.Equal(t, "owner", relations[2].Name)
	assert.Equal(t, "user:bob", relations[2].To)
	assert.Equal(t, []string{"resource:foo", "owner", "collection:admins", "child", "collection:superadmins", "member", "user:bob"}, relations[2].Via)
}

// TODO: ensure gets deleted if dependency also deleted

func TestListRelationsRec(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Create Relations

	ab := &Relation{From: "collection:a", To: "collection:b", Name: "ab"}
	err := ab.Create(ctx)
	require.NoError(t, err)

	a2b := &Relation{From: "collection:a2", To: "collection:b", Name: "a2b"}
	err = a2b.Create(ctx)
	require.NoError(t, err)

	bc := &Relation{From: "collection:b", To: "collection:c", Name: "bc"}
	err = bc.Create(ctx)
	require.NoError(t, err)

	cd := &Relation{From: "collection:c", To: "collection:d", Name: "cd"}
	err = cd.Create(ctx)
	require.NoError(t, err)

	de := &Relation{From: "collection:d", To: "collection:e", Name: "de"}
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

		assert.Equal(t, "collection:c", rels[0].From)
		assert.Equal(t, "collection:d", rels[0].To)
		assert.Equal(t, []string{"collection:c", "cd"}, rels[0].Via)

		assert.Equal(t, "collection:b", rels[1].From)
		assert.Equal(t, "collection:c", rels[1].To)
		assert.Equal(t, []string{"collection:b", "bc", "collection:c", "cd"}, rels[1].Via)

		assert.Equal(t, "collection:a2", rels[2].From)
		assert.Equal(t, "collection:b", rels[2].To)
		assert.Equal(t, []string{"collection:a2", "a2b", "collection:b", "bc", "collection:c", "cd"}, rels[2].Via)

		assert.Equal(t, "collection:a", rels[3].From)
		assert.Equal(t, "collection:b", rels[3].To)
		assert.Equal(t, []string{"collection:a", "ab", "collection:b", "bc", "collection:c", "cd"}, rels[3].Via)
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
		assert.Equal(t, "collection:c", rels[0].From)
		assert.Equal(t, "collection:d", rels[0].To)
		assert.Equal(t, []string{"cd", "collection:d"}, rels[0].Via)

		assert.Equal(t, "collection:d", rels[1].From)
		assert.Equal(t, "collection:e", rels[1].To)
		assert.Equal(t, []string{"cd", "collection:d", "de", "collection:e"}, rels[1].Via)
	})
}

func TestRelationsDelete(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	fooOwnerAdmins := newRelation("resource:foo", "owner", "collection:admins")
	err := fooOwnerAdmins.Create(ctx)
	require.NoError(t, err)

	adminsMemberAlice := newRelation("collection:admins", "member", "user:alice")
	err = adminsMemberAlice.Create(ctx)
	require.NoError(t, err)

	rels, err := Check(ctx, "resource:foo", "owner", "user:alice")
	require.NoError(t, err)

	assert.Equal(t, 1, len(rels))

	toDelete := &Relation{
		From: adminsMemberAlice.From,
		Name: adminsMemberAlice.Name,
		To:   adminsMemberAlice.To,
	}
	require.NoError(t, toDelete.Delete(ctx))

	rels, err = Check(ctx, "resource:foo", "owner", "user:alice")
	require.NoError(t, err)

	assert.Equal(t, 0, len(rels))

}
