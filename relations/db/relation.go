package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/u"
)

var Conn *pgxpool.Pool

type Relation struct {
	ID       string    `db:"_id"`
	From     EntityRef `db:"from"`
	To       EntityRef `db:"to"`
	Attrs    map[string]any
	Indirect bool
}

type EntityRef struct {
	ID   string
	Type string
}

type RelationFilter struct {
	FromID   *string `db:"from_id"`
	FromType *string `db:"from_type"`
	ToID     *string `db:"to_id"`
	ToType   *string `db:"to_type"`
	Indirect *bool
}

func (r *Relation) Create(ctx context.Context) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	if r.Attrs == nil {
		r.Attrs = map[string]any{}
	}

	query := `
	  insert into relations(_id, from_id, from_type, to_id, to_type, attrs, indirect)
	  values($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := Conn.Exec(ctx, query, r.ID, r.From.ID, r.From.Type, r.To.ID, r.To.Type, r.Attrs, r.Indirect)
	if err != nil {
		return err
	}

	return nil
}

func ListRelations(ctx context.Context, f RelationFilter) ([]Relation, error) {
	where, params := u.FilterBy(&f)
	query := `
	  select
	    _id,
	    from_id as "from.id",
	    from_type as "from.type",
	    to_id as "to.id",
	    to_type as "to.type",
	    attrs,
	    indirect
	  from relations
	` + where

	rs := []Relation{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, params...); err != nil {
		return nil, err
	}

	return rs, nil
}
