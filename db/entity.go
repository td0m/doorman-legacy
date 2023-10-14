package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var entityIDRe = regexp.MustCompile(`[a-z]+:[_a-z0-9]+`)

type Entity struct {
	ID    string
	Attrs map[string]any

	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Entity) Create(ctx context.Context) error {
	if !entityIDRe.MatchString(e.ID) {
		return errors.New("id format is not valid. must be in form of 'type:id'")
	}

	query := `
		insert into entities(id, attrs)
		values($1, $2)
		returning updated_at
	`

	if err := pgxscan.Get(ctx, pg, e, query, e.ID, e.Attrs); err != nil {
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

func (e *Entity) Update(ctx context.Context) error {
	query := `
		update entities
		set (updated_at, attrs) = (now(), $3)
		where
			id = $1 and updated_at = $2
	`

	if _, err := pg.Exec(ctx, query, e.ID, e.UpdatedAt, e.Attrs); err != nil {
		return fmt.Errorf("pg.Exec failed: %w", err)
	}

	return nil
}

func RetrieveEntity(ctx context.Context, id string) (*Entity, error) {
	query := `
		select id, attrs, updated_at
		from entities
		where id = $1
	`

	var e Entity
	if err := pgxscan.Get(ctx, pg, &e, query, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("pg failed: %w", err)
	}
	return &e, nil
}
