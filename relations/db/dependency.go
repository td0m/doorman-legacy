package db

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Dependency struct {
	RelationID   string
	DependencyID string
}

func (d *Dependency) Create(ctx context.Context) error {
	query := `
	  insert into dependencies(relation_id, dependency_id)
	  values($1, $2)
	`

	if _, err := Conn.Exec(ctx, query, d.RelationID, d.DependencyID); err != nil {
		return err
	}

	return nil
}

func ListDependencies(ctx context.Context, relationID string) ([]string, error) {
	query := `
    select dependency_id
    from dependencies
	  where relation_id = $1
  `

	deps := []string{}
	if err := pgxscan.Select(ctx, Conn, &deps, query, relationID); err != nil {
		return nil, err
	}

	return deps, nil
}
