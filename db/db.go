package db

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pg *pgxpool.Pool

func get(ctx context.Context, dest any, query string, args... any) error {
	return pgxscan.Get(ctx, pg, dest, query, args...)
}

func list(ctx context.Context, dest any, query string, args... any) error {
	return pgxscan.Select(ctx, pg, dest, query, args...)
}

func Init(ctx context.Context) error {
	var err error
	pg, err = pgxpool.New(ctx, "")
	if err != nil {
		return fmt.Errorf("pgxpool.New failed: %w", err)
	}
	return nil
}

func Close() error {
	pg.Close()
	return nil
}
