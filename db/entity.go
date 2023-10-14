package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var entityIDRe = regexp.MustCompile(`[a-z]+:[_a-z0-9]+`)

type Entity struct {
	ID    string
	Attrs map[string]any
}

func (e *Entity) Create(ctx context.Context) error {
	if !entityIDRe.MatchString(e.ID) {
		return errors.New("id format is not valid. must be in form of 'type:id'")
	}

	query := `
		insert into entities(id, attrs)
		values($1, $2)
	`

	if _, err := pg.Exec(ctx, query, e.ID, e.Attrs); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "entities_pkey" {
				return ErrAlreadyExists
			} else {
				fmt.Println(pgErr.Code)
			}
		}
		return fmt.Errorf("pg.Exec failed: %w", err)
	}

	return nil
}
