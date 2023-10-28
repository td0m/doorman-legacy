package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
)

var ErrInvalidRole = errors.New("this role does not exist for the given object type")
var ErrCycle = errors.New("cycle detected")
var ErrTupleNotFound = errors.New("tuple not found")

type Tuples struct {
	pool *pgxpool.Pool
}

func (t Tuples) Add(ctx context.Context, tuple doorman.Tuple) error {
	query := `
		insert into tuples(subject, role, object)
		values($1, $2, $3)
	`

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx failed: %w", err)
	}

	// Without locking we can get some issues with concurrent writes
	// Comment out for proof / check what will break
	if _, err := tx.Exec(ctx, `lock table tuples in access exclusive mode`); err != nil {
		return fmt.Errorf("locking table tuples failed: %w", err)
	}

	if _, err := tx.Exec(ctx, query, tuple.Subject, tuple.Role, tuple.Object); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "tuples_role_fkey" && pgErr.Code == "23503" {
				return ErrInvalidRole
			}
		}
		return err
	}

	connected, err := listConnected(ctx, tx, tuple.Object)
	if err != nil {
		return fmt.Errorf("listConnected failed: %w", err)
	}
	for _, o := range connected {
		if o == tuple.Subject {
			if err := tx.Rollback(ctx); err != nil {
				return fmt.Errorf("failed to rollback on cycle: %w", err)
			}
			return ErrCycle
		}
	}

	return tx.Commit(ctx)
}

func (t Tuples) Remove(ctx context.Context, tuple doorman.Tuple) error {
	query := `
		delete from tuples
		where (subject, role, object) = ($1, $2, $3)
	`

	tag, err := t.pool.Exec(ctx, query, tuple.Subject, tuple.Role, tuple.Object)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrTupleNotFound
	}
	return nil
}

func NewTuples(conn *pgxpool.Pool) Tuples {
	return Tuples{conn}
}

func (t Tuples) ListConnected(ctx context.Context, subject doorman.Object, inverted bool) ([]doorman.Path, error) {
	query := `
		with recursive connections as (
			select
				object, array_append(array_append('{}'::text[], role), object) as via
			from tuples
			where subject = $1

			union

			select next.object, array_append(array_append(prev.via, next.role), next.object)
			from tuples next
			inner join
				connections prev on prev.object = next.subject
			where next.object != $1
		) select via from connections
	`

	if inverted {
		query = `
		with recursive inverted_connections as (
			select
				subject, array_append(array_append('{}'::text[], role), subject) as via
			from tuples
			where object = $1

			union

			select next.subject, array_append(array_append(prev.via, next.role), next.subject)
			from tuples next
			inner join
				inverted_connections prev on prev.subject = next.object
			where next.subject != $1
		) select via from inverted_connections
	`
	}

	rows, err := t.pool.Query(ctx, query, subject)
	if err != nil {
		return nil, fmt.Errorf("exec failed: %w", err)
	}

	var paths []doorman.Path
	for rows.Next() {
		var via []string
		if err := rows.Scan(&via); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		path := make([]doorman.Connection, len(via)/2)
		for i := 0; i < len(via); i += 2 {
			path[i/2] = doorman.Connection{Role: via[i], Object: doorman.Object(via[i+1])}
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func listConnected(ctx context.Context, tx pgx.Tx, subject doorman.Object) ([]doorman.Object, error) {
	query := `
		with recursive connections as (
			select
				object
			from tuples
			where subject = $1

			union

			select next.object
			from tuples next
			inner join
				connections prev on prev.object = next.subject
			where next.object != $1
		) select object from connections
	`

	rows, err := tx.Query(ctx, query, subject)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var objects []doorman.Object
	for rows.Next() {
		o := doorman.Object("")
		if err := rows.Scan(&o); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		objects = append(objects, o)
	}

	return objects, nil
}
