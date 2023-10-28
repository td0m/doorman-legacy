package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Roles struct {
	pool *pgxpool.Pool
}

func (r Roles) Add(ctx context.Context, role doorman.Role) error {
	query := `
		insert into roles(id, verbs)
		values($1, $2)
	`

	if _, err := r.pool.Exec(ctx, query, role.ID, role.Verbs); err != nil {
		return err
	}

	return nil
}

func (r Roles) Retrieve(ctx context.Context, id string) (*doorman.Role, error) {
	query := `
		select verbs
		from roles
		where id = $1
	`

	role := doorman.Role{ID: id}

	err := r.pool.QueryRow(ctx, query, id).Scan(&role.Verbs)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &role, nil
}

func (r Roles) Update(ctx context.Context, role *doorman.Role) (error) {
	query := `
		update roles
		set verbs = $2
		where id = $1
	`

	if _, err := r.pool.Exec(ctx, query, role.ID, role.Verbs); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func NewRoles(pool *pgxpool.Pool) Roles {
	return Roles{pool}
}
