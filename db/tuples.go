package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/td0m/doorman"
	"golang.org/x/exp/slices"
)

var ErrInvalidRole = errors.New("this role does not exist for the given object type")
var ErrCycle = errors.New("cycle detected")
var ErrTupleNotFound = errors.New("tuple not found")

type Tuples struct {
	pool *pgxpool.Pool
}

func (t Tuples) Add(ctx context.Context, tuple doorman.Tuple) ([]doorman.TupleWithPath, error) {
	query := `
		insert into tuples(subject, role, object)
		values($1, $2, $3)
	`

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Without locking we can get some issues with concurrent writes
	// Comment out for proof / check what will break
	if _, err := tx.Exec(ctx, `lock table tuples in access exclusive mode`); err != nil {
		return nil, fmt.Errorf("locking table tuples failed: %w", err)
	}

	if _, err := tx.Exec(ctx, query, tuple.Subject, tuple.Role, tuple.Object); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "tuples_role_fkey" && pgErr.Code == "23503" {
				return nil, ErrInvalidRole
			}
		}
		return nil, err
	}

	connected, err := listConnectedTiny(ctx, tx, tuple.Object)
	if err != nil {
		return nil, fmt.Errorf("listConnected failed: %w", err)
	}
	for _, o := range connected {
		if o == tuple.Subject {
			if err := tx.Rollback(ctx); err != nil {
				return nil, fmt.Errorf("failed to rollback on cycle: %w", err)
			}
			return nil, ErrCycle
		}
	}

	newTuples, err := t.tuplesFrom(ctx, tx, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuplesFrom failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return newTuples, nil
}

func (t Tuples) Remove(ctx context.Context, tuple doorman.Tuple) ([]doorman.TupleWithPath, error) {
	query := `
		delete from tuples
		where (subject, role, object) = ($1, $2, $3)
	`

	tag, err := t.pool.Exec(ctx, query, tuple.Subject, tuple.Role, tuple.Object)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() == 0 {
		return nil, ErrTupleNotFound
	}

	tuples, err := t.tuplesFrom(ctx, t.pool, tuple)
	if err != nil {
		return nil, fmt.Errorf("tuplesFrom failed: %w", err)
	}

	return tuples, nil
}

func NewTuples(conn *pgxpool.Pool) Tuples {
	return Tuples{conn}
}

func (t Tuples) ListTuplesForRole(ctx context.Context, role string) ([]doorman.Tuple, error) {
	query := `
		select subject, object
		from tuples
		where role = $1
	`

	rows, err := t.pool.Query(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	var tuples []doorman.Tuple
	for rows.Next() {
		t := doorman.Tuple{Role: role}
		if err := rows.Scan(&t.Subject, &t.Object); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		tuples = append(tuples, t)
	}
	return tuples, nil
}

func (t Tuples) ListConnected(ctx context.Context, subject doorman.Object, inverted bool) ([]doorman.Path, error) {
	return listConnected(ctx, t.pool, subject, inverted)
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func listConnected(ctx context.Context, tx querier, subject doorman.Object, inverted bool) ([]doorman.Path, error) {
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

	rows, err := tx.Query(ctx, query, subject)
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

func listConnectedTiny(ctx context.Context, tx pgx.Tx, subject doorman.Object) ([]doorman.Object, error) {
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

func (ts *Tuples) tuplesFrom(ctx context.Context, tx querier, tuple doorman.Tuple) ([]doorman.TupleWithPath, error) {
	newTuples := []doorman.TupleWithPath{}
	{
		tupleChildren := []doorman.Path{{}}
		if tuple.Object.Type() == "group" {
			connections, err := listConnected(ctx, tx, tuple.Object, false)
			if err != nil {
				return nil, fmt.Errorf("tuples.ListConnected(subj, false) failed: %w", err)
			}
			tupleChildren = append(tupleChildren, connections...)
		}

		tupleParents := []doorman.Path{{}}
		if tuple.Subject.Type() == "group" {
			connections, err := listConnected(ctx, tx, tuple.Subject, true)
			if err != nil {
				return nil, fmt.Errorf("tuples.ListConnected(obj, true) failed: %w", err)
			}
			tupleParents = append(tupleParents, connections...)
		}

		for _, child := range tupleChildren {
			for _, parent := range tupleParents {
				t := doorman.TupleWithPath{
					Tuple: doorman.Tuple{
						Subject: tuple.Subject,
						Role:    tuple.Role,
						Object:  tuple.Object,
					},
					Path: doorman.Path{},
				}
				if len(parent) > 0 {
					t.Subject = parent[len(parent)-1].Object
					path := parent[:len(parent)-1]
					slices.Reverse(path)
					t.Path = append(t.Path, path...)
				}
				if len(child) > 0 {
					t.Object = child[len(child)-1].Object
					t.Role = child[len(child)-1].Role
					t.Path = append(t.Path, child[:len(child)-1]...)
				}

				if throughGroupsOnly(t.Path) {
					newTuples = append(newTuples, t)
				}
			}
		}
	}

	return newTuples, nil
}

// e.g. user:alice -> item:1 -> item:2 should not connect user:alice with item:2
// why? because it's not a group
func throughGroupsOnly(path doorman.Path) bool {
	for _, conn := range path {
		if conn.Object.Type() != "group" {
			return false
		}
	}

	return true
}
