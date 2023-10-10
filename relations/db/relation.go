package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/u"
)

var Conn *pgxpool.Pool

type Relation struct {
	Cache bool // Decided whether we use the unlogged cache table or source of truth

	ID   string    `db:"_id"`
	From EntityRef `db:"from"`
	To   EntityRef `db:"to"`
	Name *string
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
	FromID          *string `db:"from_id"`
	FromType        *string `db:"from_type"`
	ToID            *string `db:"to_id"`
	ToType          *string `db:"to_type"`
	PaginationToken *string `db:"_id" op:">"`
	Name            *string
}

func (r *Relation) Delete(ctx context.Context) error {
	query := `delete from ` + r.table() + ` where _id=$1`
	_, err := Conn.Exec(ctx, query, r.ID)
	return err
}

func (r *Relation) Create(ctx context.Context) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	if r.Name != nil && len(*r.Name) == 0 {
		r.Name = nil
	}

	query := `
	  insert into ` + r.table() + `(_id, from_id, from_type, to_id, to_type, name)
	  values($1, $2, $3, $4, $5, $6)
	`

	_, err := Conn.Exec(ctx, query, r.ID, r.From.ID, r.From.Type, r.To.ID, r.To.Type, r.Name)
	if err != nil {
		return err
	}

	return nil
}

// We always get from cache
func Get(ctx context.Context, id string) (*Relation, error) {
	query := `
	  select
	    _id,
	    name,
	    from_id as "from.id",
	    from_type as "from.type",
	    to_id as "to.id",
	    to_type as "to.type"
	  from cache
	  where
	    _id = $1
	`

	var r Relation
	if err := pgxscan.Get(ctx, Conn, &r, query, id); err != nil {
		return nil, err
	}
	return &r, nil
}

func ListLinkedToWithDependencies(ctx context.Context, e EntityRef) ([]RelationWithDeps, error) {
	query := `
	  select
	    _id,
	    name,
	    from_id as "from.id",
	    from_type as "from.type",
	    to_id as "to.id",
	    to_type as "to.type",
	    array(select relation_id from dependencies where cache_id = _id)::text[] as dependency_ids
	  from cache
	  where
	    to_type = $1 and to_id = $2
	`

	rs := []RelationWithDeps{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, e.Type, e.ID); err != nil {
		return nil, err
	}
	return rs, nil
}

func (r Relation) table() string {
	return tableName(r.Cache)
}

func ListLinkedFromWithDependencies(ctx context.Context, e EntityRef) ([]RelationWithDeps, error) {
	query := `
	  select
	    _id,
	    name,
	    from_id as "from.id",
	    from_type as "from.type",
	    to_id as "to.id",
	    to_type as "to.type",
	    array(select relation_id from dependencies where cache_id = _id)::text[] as dependency_ids
	  from cache
	  where
	    from_type = $1 and from_id = $2
	`

	rs := []RelationWithDeps{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, e.Type, e.ID); err != nil {
		return nil, err
	}
	return rs, nil
}

func ListRelations(ctx context.Context, f RelationFilter, cache bool) ([]Relation, error) {
	where, params := u.FilterBy(&f)
	query := `
	  select
	    _id,
	    name,
	    from_id as "from.id",
	    from_type as "from.type",
	    to_id as "to.id",
	    to_type as "to.type"
	  from ` + tableName(cache) + ` ` + where + ` limit 1000`

	rs := []Relation{}
	if err := pgxscan.Select(ctx, Conn, &rs, query, params...); err != nil {
		return nil, err
	}

	return rs, nil
}

func tableName(cache bool) string {
	if cache {
		return "cache"
	}
	return "relations"
}
