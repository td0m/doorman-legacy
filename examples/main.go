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

	dom, err := entities.Create(ctx, entities.CreateRequest{Type: "user"})
	check(err)

	admins, err := entities.Create(ctx, entities.CreateRequest{Type: "collection"})
	check(err)

	post1, err := entities.Create(ctx, entities.CreateRequest{
		ID:   "posts/" + xid.New().String(),
		Type: "resource",
	})
	check(err)

	_, err = relations.Create(ctx, relations.CreateRequest{
		From: relations.Entity{ID: dom.ID, Type: dom.Type},
		To: relations.Entity{ID: admins.ID, Type: admins.Type},
	})

	check(err)
	_, err = relations.Create(ctx, relations.CreateRequest{
		From: relations.Entity{ID: admins.ID, Type: admins.Type},
		To:   relations.Entity{ID: post1.ID, Type: post1.Type},
	})
	check(err)

	rs, err := relations.List(ctx, relations.ListRequest{
		From: &relations.Entity{ID: dom.ID, Type: dom.Type},
		To:   &relations.Entity{ID: post1.ID, Type: post1.Type},
	})
	check(err)

	fmt.Println("relations", len(rs))
	for _, r := range rs {
		fmt.Println(r.ID)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
