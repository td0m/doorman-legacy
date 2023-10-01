package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/u"
)

var Conn *pgxpool.Pool

type Relation struct {
	ID    string `db:"_id"`
	From  string
	To    string
	Attrs map[string]any
}

type RelationFilter struct {
	From *string `db:"\"from\""`
	To   *string `db:"\"to\""`
}

func (r *Relation) Create(ctx context.Context) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	if r.Attrs == nil {
		r.Attrs = map[string]any{}
	}

	query := `
	  insert into relations(_id, "from", "to", attrs)
	  values($1, $2, $3, $4)
	`

	fmt.Printf("%+v", r)

	_, err := Conn.Exec(ctx, query, r.ID, r.From, r.To, r.Attrs)
	if err != nil {
		return err
	}

	return nil
}

func ListRelations(ctx context.Context, f RelationFilter) ([]Relation, error) {
	where, params := u.FilterBy(&f)
	query := `
	  select _id, "from", "to", attrs
	  from relations
	` + where

	rs := []Relation{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, params...); err != nil {
		return nil, err
	}

	return rs, nil
}
