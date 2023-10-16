package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pg *pgxpool.Pool

func Init(ctx context.Context) error {
	var err error
	pg, err = pgxpool.New(ctx, "")
	if err != nil {
		return fmt.Errorf("pgxpool.New failed: %w", err)
	}

	if err := pg.Ping(ctx); err != nil {
		return fmt.Errorf("pg.Ping failed: %w", err)
	}
	return nil
}

func Close() error {
	pg.Close()
	return nil
}
