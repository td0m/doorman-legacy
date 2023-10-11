package db

import (
	"context"
	"net/http"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/td0m/poc-doorman/errs"
)

type Type struct {
	ID    string `db:"_id"`
	Attrs map[string]any

	UpdatedAt time.Time
}

func (t *Type) Create(ctx context.Context) error {
	if !wordRe.MatchString(t.ID) {
		return errs.New(http.StatusBadRequest, "ID is invalid, must be an alphanumeric word starting with a letter.")
	}

	if t.Attrs == nil {
		t.Attrs = map[string]any{}
	}

	query := `
	  insert into entity_types(_id, attrs)
	  values($1, $2)
	`

	_, err := Conn.Exec(ctx, query, t.ID, t.Attrs)
	return err
}

func (t *Type) Update(ctx context.Context) error {
	query := `
	  update entity_types
	  set
	    updated_at = now(),
	    attrs = $3
	  where
	    _id = $1 and
	    updated_at = $2
	  returning updated_at
  `

	err := pgxscan.Get(ctx, Conn, t, query, t.ID, t.UpdatedAt, t.Attrs)
	return err
}

func RetrieveType(ctx context.Context, id string) (*Type, error) {
	query := `
	  select _id, attrs, updated_at
	  from entity_types
	  where _id = $1
	`

	var t Type
	err := pgxscan.Get(ctx, Conn, &t, query, id)
	return &t, err
}

func ListTypes(ctx context.Context) ([]Type, error) {
	query := `
	  select _id, attrs, updated_at
	  from entity_types
	`

	var ts []Type
	err := pgxscan.Get(ctx, Conn, &ts, query)

	return ts, err
}
