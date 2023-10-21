package server

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
)

func cleanup(t *testing.T) {
	ctx := context.Background()
	_, err := db.Conn().Exec(ctx, `delete from dependencies;`)
	if err != nil {
		t.Fail()
	}
	_, err = db.Conn().Exec(ctx, `delete from relations;`)
	if err != nil {
		t.Fail()
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

	if err := db.Init(ctx); err != nil {
		panic(err)
	}

	m.Run()
}
func TestCheckStored(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	server := NewDoormanServer()

	in := &pb.CheckRequest{
		User:   "user:alice",
		Name:   "owner",
		Object: "product:apple",
	}

	// Check before
	res, err := server.Check(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Make relation
	relation := &db.Relation{
		From: "product:apple",
		Name: "owner",
		To:   "user:alice",
	}
	err = relation.Create(ctx)
	require.NoError(t, err)

	// Check after
	res, err = server.Check(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)
}

func TestCheckComputed(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	server := NewDoormanServer()

	in := &pb.CheckRequest{
		User:   "user:alice",
		Name:   "owner",
		Object: "product:apple",
	}

	// Check before
	res, err := server.Check(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)

	// Make relations
	appleOnShelf1 := &db.Relation{
		From: "product:apple",
		Name: "on",
		To:   "shelf:1",
	}
	err = appleOnShelf1.Create(ctx)
	require.NoError(t, err)

	shelf1Owner := &db.Relation{
		From: "shelf:1",
		Name: "owner",
		To:   "user:alice",
	}
	err = shelf1Owner.Create(ctx)
	require.NoError(t, err)

	// Check after
	res, err = server.Check(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, true, res.Connected)

	in.User = "user:randomusernoperms"
	res, err = server.Check(ctx, in)
	assert.NoError(t, err)
	assert.Equal(t, false, res.Connected)
}
