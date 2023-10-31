package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Relations struct {
	conn querier
}

func (r Relations) WithTx(tx pgx.Tx) *Relations {
	return &Relations{conn: tx}
}

// 1m users, 1k posts
//  -> new post = 1 million new ACLs... that is bad
//  but we don't need to give them direct access... we do it smart
//  users access feed, feed access rungways

func (rs Relations) Check(ctx context.Context, r doorman.Relation) (bool, error) {
	query := `
		select 1
		from relations
		where (subject, verb, object) = ($1, $2, $3)
		limit 1
	`

	rows, err := rs.conn.Query(ctx, query, r.Subject, r.Verb, r.Object)
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
		insert into relations(subject, verb, object)
		values($1, $2, $3)
		on conflict do nothing
	`

	if _, err := rs.conn.Exec(ctx, query, r.Subject, r.Verb, r.Object); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (rs Relations) Remove(ctx context.Context, r doorman.Relation) error {
	query := `
		delete from relations
		where (subject, verb, object) = ($1, $2, $3)
	`
	if _, err := rs.conn.Exec(ctx, query, r.Subject, r.Verb, r.Object); err != nil {
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
	return Relations{conn: pool}
}
