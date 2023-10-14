package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/u"
)

type Cache struct {
	ID   string
	Name *string
	From string
	To   string

	Via []string // TODO: remove
}

type RelationFilter struct {
	AfterID  *string `db:"id" op:">"`
	From     *string `db:"\"from\""`
	FromType *string `db:"from_type"`
	Name     *string
	To       *string `db:"\"to\""`
	ToType   *string `db:"to_type"`
}

type Executioner interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

func depsToID(deps []string) string {
	return "cache: " + strings.Join(deps, " ")
}

func (c *Cache) Delete(ctx context.Context) error {
	query := `delete from cache where id = $1`
	if _, err := pg.Exec(ctx, query, c.ID); err != nil {
		return fmt.Errorf("pg.Exec failed: %w", err)
	}
	return nil
}

func (c *Cache) Create(ctx context.Context, conn Executioner) error {
	if c.ID == "" {
		if len(c.Via) > 0 {
			c.ID = depsToID(c.Via)
		} else {
			c.ID = xid.New().String()
		}
	}

	query := `
		insert into cache(id, "from", "to", name)
		values($1, $2, $3, $4)
		on conflict on constraint cache_pkey
		do nothing
	`

	if _, err := conn.Exec(ctx, query, c.ID, c.From, c.To, c.Name); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}

	return nil
}

func ListRelationsOrCache(ctx context.Context, table string, f RelationFilter) ([]Cache, error) {
	cols := `id, name, "from", "to"`

	where, params := u.FilterBy(&f)
	query := `
		select ` + cols + `
		from ` + table + `
		` + where + `
		limit 1000
	`

	var caches []Cache

	if err := pgxscan.Select(ctx, pg, &caches, query, params...); err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return caches, nil
}
