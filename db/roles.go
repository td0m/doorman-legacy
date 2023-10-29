package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

type Roles struct {
	conn querier
}

func (r Roles) WithTx(tx pgx.Tx) *Roles {
	return &Roles{conn: tx}
}

func (r Roles) Add(ctx context.Context, role doorman.Role) error {
	query := `
		insert into roles(id, verbs)
		values($1, $2)
	`

	if _, err := r.conn.Exec(ctx, query, role.ID, role.Verbs); err != nil {
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

	err := r.conn.QueryRow(ctx, query, id).Scan(&role.Verbs)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrInvalidRole
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return &role, nil
}

func (r Roles) Remove(ctx context.Context, id string) error {
	query := `
		delete from roles where id=$1
	`

	if _, err := r.conn.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (r Roles) Upsert(ctx context.Context, role *doorman.Role) error {
	query := `
		insert into roles(id, verbs)
		values($1, $2)
		on conflict(id) do update
			set verbs = $2
	`

	if _, err := r.conn.Exec(ctx, query, role.ID, role.Verbs); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func NewRoles(pool *pgxpool.Pool) Roles {
	return Roles{pool}
}
