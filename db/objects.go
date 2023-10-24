package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Objects struct {
	pool *pgxpool.Pool
}

func (o Objects) Add(ctx context.Context, obj doorman.Object) error {
	query := `
		insert into objects(id)
		values($1)
	`

	if _, err := o.pool.Exec(ctx, query, obj); err != nil {
		return err
	}

	return nil
}

func NewObjects(pool *pgxpool.Pool) Objects {
	return Objects{pool}
}
