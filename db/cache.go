package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
)

type Cache struct {
	ID   string
	Name *string
	From string
	To   string

	Via []string
}

type CacheFilter struct {
	From     *string
	FromType *string `db:"from_type"`
	To       *string
	ToType   *string `db:"to_type"`
	Name     *string
}

func (c *Cache) Create(ctx context.Context, tx pgx.Tx) error {
	if c.ID == "" {
		if len(c.Via) > 0 {
			c.ID = strings.Join(c.Via, " ")
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

	if _, err := tx.Exec(ctx, query, c.ID, c.From, c.To, c.Name); err != nil {
		return fmt.Errorf("tx.Exec failed: %w", err)
	}

	return nil
}

func ListCaches(ctx context.Context, f CacheFilter) ([]RecRelation, error) {
}
