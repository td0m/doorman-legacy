package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

type Entity struct {
	ID    string
	Attrs map[string]any
}

func (e *Entity) Create(ctx context.Context) error {
	if e.Attrs == nil {
		e.Attrs = map[string]any{}
	}

	query := `
	  insert into entities(_id, attrs)
	  values($1, $2)
	`
	fmt.Printf("%+v", e)

	_, err := Conn.Exec(ctx, query, e.ID, e.Attrs)
	if err != nil {
		return err
	}

	return nil
}
