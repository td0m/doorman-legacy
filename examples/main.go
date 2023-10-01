package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/entities"
	entitiesdb "github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/relations"
	relationsdb "github.com/td0m/poc-doorman/relations/db"
)

func main() {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	check(err)
	defer pgDoorman.Close()

	entitiesdb.Conn = pgDoorman
	relationsdb.Conn = pgDoorman

	dom, err := entities.Create(ctx, entities.Entity{
		Type: "user",
	})
	check(err)

	fmt.Println(dom)

	post1, err := entities.Create(ctx, entities.Entity{
		ID: "posts/" + xid.New().String(),
		Type: "resource",
	})
	check(err)

	fmt.Println(post1)

	r, err := relations.Create(ctx, relations.CreateRequest{
		From: dom.ID,
		To: post1.ID,
	})
	check(err)
	fmt.Println(r)

	rs, err := relations.List(ctx, relations.ListRequest{
		From: dom.ID,
		To: post1.ID,
	})
	check(err)
	fmt.Println("relations", rs)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
