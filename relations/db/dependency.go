package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Dependency struct {
	CacheID    string
	RelationID string
}

func (d *Dependency) Create(ctx context.Context) error {
	query := `
	  insert into dependencies(cache_id, relation_id)
	  values($1, $2)
	`

	if _, err := Conn.Exec(ctx, query, d.CacheID, d.RelationID); err != nil {
		return err
	}

	return nil
}

func ListDependencies(ctx context.Context, cacheID string) ([]string, error) {
	query := `
    select relation_id
    from dependencies
	  where cache_id = $1
  `

	var relationIDs []string
	if err := pgxscan.Select(ctx, Conn, &relationIDs, query, cacheID); err != nil {
		return nil, err
	}

	return relationIDs, nil
}
