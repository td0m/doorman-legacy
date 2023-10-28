package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Changes struct {
	pool *pgxpool.Pool
}

type ChangeFilter struct {
	PaginationToken *string
}

func (cs Changes) Add(ctx context.Context, c doorman.Change) error {
	return nil
}

func (cs Changes) List(ctx context.Context, f ChangeFilter) ([]doorman.Change, error) {
	return nil, nil
}

func NewChanges(pool *pgxpool.Pool) Changes {
	return Changes{pool}
}
