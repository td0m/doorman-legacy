package store

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Postgres struct {
	conn *pgxpool.Pool
}

func NewPostgres(conn *pgxpool.Pool) Postgres {
	return Postgres{conn: conn}
}

func (p Postgres) Add(ctx context.Context, t Tuple) error {
	query := `
		insert into tuples(u, label, v)
		values($1, $2, $3)
	`

	if _, err := p.conn.Exec(ctx, query, t.U, t.Label, t.V); err != nil {
		return err
	}
	return nil
}

func (p Postgres) Remove(ctx context.Context, id string) error {
	return nil
}

func (p Postgres) Check(ctx context.Context, s doorman.Set, e doorman.Element) (bool, error) {
	query := `
		select v
		from tuples
		where
			u      = $1 and
			label  = $2 and
			v      = $3
	`

	var items []doorman.Element
	if err := pgxscan.Select(ctx, p.conn, &items, query, s.U, s.Label, e); err != nil {
		return false, fmt.Errorf("select failed: %w", err)
	}

	return len(items) > 0, nil
}

func (p Postgres) ListElements(ctx context.Context, set doorman.Set) ([]doorman.Element, error) {
	query := `
		select v
		from tuples
		where
			u      = $1 and
			label  = $2
	`

	var items []doorman.Element
	rows, err := p.conn.Query(ctx, query, set.U, set.Label)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	for rows.Next() {
		var item doorman.Element
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
