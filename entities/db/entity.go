package db

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/errs"
)

var wordRe = regexp.MustCompile(`^[a-z]+$`)

var Conn *pgxpool.Pool

type Entity struct {
	ID    string `db:"_id"`
	Type  string `db:"_type"`
	Attrs map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Entity) Create(ctx context.Context) error {
	if e.ID == "" {
		e.ID = xid.New().String()
	} else if strings.Contains(e.ID, " ") {
		return errs.New(http.StatusBadRequest, "ID cannot contain spaces")
	}
	if len(e.Attrs) == 0 {
		e.Attrs = nil
	}
	if !wordRe.MatchString(e.Type) {
		return errs.New(http.StatusBadRequest, "Type is invalid, must be an all lowercase word.")
	}

	query := `
	  insert into entities(_id, _type, attrs)
	  values($1, $2, $3)
	`

	_, err := Conn.Exec(ctx, query, e.ID, e.Type, e.Attrs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Entity) Update(ctx context.Context) error {
	query := `
	  update entities
	  set
	    updated_at = now(),
	    attrs = $4
	  where
	    _id = $1 and
	    _type = $2 and
	    updated_at = $3
	  returning updated_at
	`

	err := pgxscan.Get(ctx, Conn, e, query, e.ID, e.Type, e.UpdatedAt, e.Attrs)
	if err != nil {
		return err
	}

	return nil
}

func (e *Entity) Delete(ctx context.Context) error {
	query := `
	  delete from entities
	  where
	    _id = $1 and
	    _type = $2
	`

	_, err := Conn.Exec(ctx, query, e.ID, e.Type)
	if err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context, id string, typ string) (*Entity, error) {
	query := `
	  select
	    _id,
	    _type,
	    attrs,
			created_at
			updated_at
	  from entities
	  where
	    _id = $1 and
	    _type = $2
	`
	var e Entity
	if err := pgxscan.Get(ctx, Conn, &e, query, id, typ); err != nil {
		return nil, err
	}

	return &e, nil
}
