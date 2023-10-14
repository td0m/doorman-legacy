package db

import (
	"context"
	"fmt"

	"github.com/rs/xid"
)


type Entity struct {
	ID    string
	Attrs map[string]any
}

func (e *Entity) Create(ctx context.Context) error {
	if e.ID == "" {
		e.ID = xid.New().String()
	}

	query := `
		insert into entities(id, attrs)
		values($1, $2)
	`

	if _, err := pg.Exec(ctx, query, e.ID, e.Attrs); err != nil {
		return fmt.Errorf("pg.Exec failed: %w", err)
	}

	return nil
}
