package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

type Entity struct {
	ID    string
	Type  string
	Attrs map[string]any
}

func (e *Entity) Create(ctx context.Context) error {
	if e.Attrs == nil {
		e.Attrs = map[string]any{}
	}

	query := `
	  insert into entities(_id, _type, attrs)
	  values($1, $2, $3)
	`

	_, err := Conn.Exec(ctx, query, e.ID, e.Type, e.Attrs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Entity) Delete(ctx context.Context) error {
	query := `
	  delete from entities
	  where
	    _id = $1 and
	    _type = $2
	`

	_, err := Conn.Exec(ctx, query, e.ID, e.Type)
	if err != nil {
		return err
	}

	return nil
}

