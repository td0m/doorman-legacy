package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/entities"
	"github.com/td0m/poc-doorman/ui"
	entitiesdb "github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/relations"
	relationsdb "github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
)

func main() {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)
	defer pgDoorman.Close()

	entitiesdb.Conn = pgDoorman
	relationsdb.Conn = pgDoorman

	dom, err := entities.Create(ctx, entities.CreateRequest{Type: "user"})
	u.Check(err)

	admins, err := entities.Create(ctx, entities.CreateRequest{Type: "collection"})
	u.Check(err)

	post1, err := entities.Create(ctx, entities.CreateRequest{
		ID:   xid.New().String(),
		Type: "post",
	})
	u.Check(err)

	_, err = relations.Create(ctx, relations.CreateRequest{
		From: relations.Entity{ID: dom.ID, Type: dom.Type},
		To: relations.Entity{ID: admins.ID, Type: admins.Type},
	})

	u.Check(err)
	_, err = relations.Create(ctx, relations.CreateRequest{
		From: relations.Entity{ID: admins.ID, Type: admins.Type},
		To:   relations.Entity{ID: post1.ID, Type: post1.Type},
	})
	u.Check(err)

	rs, err := relations.List(ctx, relations.ListRequest{
		From: &relations.Entity{ID: dom.ID, Type: dom.Type},
		To:   &relations.Entity{ID: post1.ID, Type: post1.Type},
	})
	u.Check(err)

	fmt.Println("relations", len(rs.Data))
	for _, r := range rs.Data {
		fmt.Println(r.ID)
	}

	u.Check(ui.Serve())
}

