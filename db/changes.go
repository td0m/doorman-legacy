package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Changes struct {
	conn querier
}

func (c Changes) WithTx(tx pgx.Tx) *Changes {
	return &Changes{conn: tx}
}

type ChangeFilter struct {
	PaginationToken *string `db:"id" op:">"`
}

func (cs Changes) Add(ctx context.Context, c doorman.Change) error {
	query := `
		insert into changes(id, type, payload) values($1, $2, $3)
	`

	if _, err := cs.conn.Exec(ctx, query, c.ID, c.Type, []byte(c.Payload)); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

func (cs Changes) List(ctx context.Context, f ChangeFilter) ([]doorman.Change, error) {
	where, params := filterBy(&f)

	query := `
		select id, type, payload, created_at
		from changes
	` + where + `
		order by id
	`

	fmt.Println("where", where, f)

	rows, err := cs.conn.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	changes := []doorman.Change{}
	for rows.Next() {
		change := doorman.Change{}
		if err := rows.Scan(&change.ID, &change.Type, &change.Payload, &change.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		changes = append(changes, change)
	}

	return changes, nil
}

func NewChanges(pool *pgxpool.Pool) Changes {
	return Changes{pool}
}
