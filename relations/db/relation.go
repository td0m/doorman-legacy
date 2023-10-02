package db

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/u"
)

var Conn *pgxpool.Pool

type Relation struct {
	ID       string    `db:"_id"`
	From     EntityRef `db:"from"`
	To       EntityRef `db:"to"`
	Attrs    map[string]any
	Indirect bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RelationWithDeps struct {
	Relation
	DependencyIDs []string `db:"dependency_ids"`
}

type EntityRef struct {
	ID   string
	Type string
}

type RelationFilter struct {
	FromID   *string `db:"from_id"`
	FromType *string `db:"from_type"`
	ToID     *string `db:"to_id"`
	ToType   *string `db:"to_type"`
	Indirect *bool
}

func Get(ctx context.Context, id string) (*Relation, error) {
	query := `
	  select
	    _id,
	    attrs,
	    created_at,
	    from_id as "from.id",
	    from_type as "from.type",
	    indirect,
	    to_id as "to.id",
	    to_type as "to.type",
	    updated_at
	  from relations
	  where
	    _id = $1
	`

	var r Relation
	if err := pgxscan.Get(ctx, Conn, &r, query, id); err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *Relation) Update(ctx context.Context) error {
	query := `
	  update relations
	  set
	    updated_at = now(),
	    attrs = $3
	  where
	    _id = $1 and
	    updated_at = $2
	  returning updated_at
	`

	err := pgxscan.Get(ctx, Conn, r, query, r.ID, r.UpdatedAt, r.Attrs)
	if err != nil {
		return err
	}

	return nil
}
func (r *Relation) Delete(ctx context.Context) error {
	query := `delete from relations where _id=$1`
	_, err := Conn.Exec(ctx, query, r.ID)
	return err
}

func (r *Relation) Create(ctx context.Context) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	if len(r.Attrs) == 0 {
		r.Attrs = nil
	}

	query := `
	  insert into relations(_id, from_id, from_type, to_id, to_type, attrs, indirect)
	  values($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := Conn.Exec(ctx, query, r.ID, r.From.ID, r.From.Type, r.To.ID, r.To.Type, r.Attrs, r.Indirect)
	if err != nil {
		return err
	}

	return nil
}

func ListLinkedToWithDependencies(ctx context.Context, e EntityRef) ([]RelationWithDeps, error) {
	query := `
	  select
	    _id,
	    attrs,
	    created_at,
	    from_id as "from.id",
	    from_type as "from.type",
	    indirect,
	    to_id as "to.id",
	    to_type as "to.type",
	    updated_at,
	    array(select dependency_id from dependencies where relation_id = _id)::text[] as dependency_ids
	  from relations
	  where
	    to_type = $1 and to_id = $2
	`

	rs := []RelationWithDeps{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, e.Type, e.ID); err != nil {
		return nil, err
	}
	return rs, nil
}

func ListLinkedFromWithDependencies(ctx context.Context, e EntityRef) ([]RelationWithDeps, error) {
	query := `
	  select
	    _id,
	    attrs,
	    created_at,
	    from_id as "from.id",
	    from_type as "from.type",
	    indirect,
	    to_id as "to.id",
	    to_type as "to.type",
	    updated_at,
	    array(select dependency_id from dependencies where relation_id = _id)::text[] as dependency_ids
	  from relations
	  where
	    from_type = $1 and from_id = $2
	`

	rs := []RelationWithDeps{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, e.Type, e.ID); err != nil {
		return nil, err
	}
	return rs, nil
}

func ListRelations(ctx context.Context, f RelationFilter) ([]Relation, error) {
	where, params := u.FilterBy(&f)
	query := `
	  select
	    _id,
	    attrs,
	    created_at,
	    from_id as "from.id",
	    from_type as "from.type",
	    indirect,
	    to_id as "to.id",
	    to_type as "to.type",
	    updated_at
	  from relations
	` + where

	rs := []Relation{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, params...); err != nil {
		return nil, err
	}

	return rs, nil
}
