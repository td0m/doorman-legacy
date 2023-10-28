package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Relations struct {
	pool *pgxpool.Pool
}

// 1m users, 1k posts
//  -> new post = 1 million new ACLs... that is bad
//  but we don't need to give them direct access... we do it smart
//  users access feed, feed access rungways

func (rs Relations) Check(ctx context.Context, r doorman.Relation) (bool, error) {
	query := `
		select path
		from relations
		where (subject, verb, object) = ($1, $2, $3)
		limit 1
	`

	rows, err := rs.pool.Query(ctx, query, r.Subject, r.Verb, r.Object)
	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	found := false
	for rows.Next() {
		_ = rows.Scan()
		found = true
	}

	return found, nil
}

func (rs Relations) Add(ctx context.Context, r doorman.Relation) error {
	// fmt.Println("add", r.String())
	// why on conflict? remove it and run the parallel test
	// alternative? use a locking tx and read tuples then...
	query := `
		insert into relations(subject, verb, object, path)
		values($1, $2, $3, $4)
		on conflict do nothing
	`

	if _, err := rs.pool.Exec(ctx, query, r.Subject, r.Verb, r.Object, r.Path); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (rs Relations) Remove(ctx context.Context, r doorman.Relation) error {
	query := `
		delete from relations
		where (subject, verb, object, path) = ($1, $2, $3, $4)
	`
	if _, err := rs.pool.Exec(ctx, query, r.Subject, r.Verb, r.Object, r.Path); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

type RelationFilter struct {
	Subject *string
	Verb    *string
}

func (rs Relations) List(ctx context.Context, f RelationFilter) ([]doorman.Relation, error) {
	panic("f")
}

func NewRelations(pool *pgxpool.Pool) Relations {
	return Relations{pool: pool}
}
