package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/xid"
	// "golang.org/x/exp/slog"
)

var ErrCycle = errors.New("cycle detected")
var ErrFkeyFrom = errors.New("entity 'from' not found")
var ErrFkeyTo = errors.New("entity 'to' not found")

type RecRelation struct {
	Via []string
	Relation
}

type Relation struct {
	ID   string
	Name *string
	From string
	To   string
	// TODO: consider attrs,
}

func extractType(id string) string {
	return strings.Split(id, ":")[0]
}

func (r Relation) Validate() error {
	if r.From == r.To {
		return fmt.Errorf("connecting to self not allowed")
	}

	return nil
}

func (r *Relation) Create(ctx context.Context) error {
	if err := r.Validate(); err != nil {
		return err
	}

	if r.ID == "" {
		r.ID = xid.New().String()
	}

	tx, err := pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pg.Begin failed: %w", err)
	}

	// Rolls back if not committed properly
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// slog.Error("tx.Rollback failed in relation.Create: ", err)
		}
	}()

	// Concurrent adding and removing relations can lead to dangling cache without this
	// Remove it and run tests for proof it's needed
	if _, err := tx.Exec(ctx, `lock table relations in access exclusive mode`); err != nil {
		return fmt.Errorf("locking table relations failed: %w", err)
	}

	query := `
		insert into relations(id, "from", "to", name)
		values($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, query, r.ID, r.From, r.To, r.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" && pgErr.ConstraintName == "relations.fkey-from" {
				return ErrFkeyFrom
			} else if pgErr.Code == "23503" && pgErr.ConstraintName == "relations.fkey-to" {
				return ErrFkeyTo
			}
		}
		return fmt.Errorf("tx.Exec failed to insert relation: %w", err)
	}

	caches, err := listDerivativeCaches(ctx, tx, *r)
	if err != nil {
		return fmt.Errorf("listDerivativeCaches failed: %w", err)
	}

	for _, cache := range caches {
		if err := cache.Create(ctx, tx); err != nil {
			return fmt.Errorf("cache.Create failed: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit failed: %w", err)
	}

	return nil
}

func listDerivativeCaches(ctx context.Context, tx pgxscan.Querier, r Relation) ([]Cache, error) {
	froms, err := listRecRelationsTo(ctx, tx, r.From)
	if err != nil {
		return nil, fmt.Errorf("listRecRelationsTo failed: %w", err)
	}
	tos, err := listRecRelationsFrom(ctx, tx, r.To)
	if err != nil {
		return nil, fmt.Errorf("listRecRelationsFrom failed: %w", err)
	}

	// Because of the nature of cycles, this will always match.
	// No need for the second statement.
	for _, from := range froms {
		if from.From == r.To {
			return nil, ErrCycle
		}
	}
	// for _, to := range to {
	// 	if to.To == r.From {
	// 		return ErrCycle
	// 	}
	// }

	froms = append(froms, RecRelation{Relation: r})
	tos = append(tos, RecRelation{Relation: r})

	caches := make([]Cache, len(froms)*len(tos))
	for i, from := range froms {
		for j, to := range tos {
			deps := append(append(from.Via, r.ID), to.Via...)
			caches[i+j*len(froms)] = Cache{
				ID:   depsToID(deps),
				From: from.From,
				To:   to.To,
				Name: nil,
			}
		}
	}

	return caches, nil
}

func listRecRelationsTo(ctx context.Context, tx pgxscan.Querier, id string) ([]RecRelation, error) {
	query := `
		with recursive relate_to as(
			select
				id,
				"from",
				"to",
				array_append('{}'::text[], id) as via
			from relations
			where "to" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				array_append(relate_to.via, r.id) as via
			from relations r
			inner join relate_to on relate_to."from" = r."to"
			where r."from" != $1
		) select * from relate_to order by id desc
	`

	var relations []RecRelation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func listRecRelationsFrom(ctx context.Context, tx pgxscan.Querier, id string) ([]RecRelation, error) {
	query := `
		with recursive relate_from as(
			select
				id,
				"from",
				"to",
				array_append('{}'::text[], id) as via
			from relations
			where "from" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				array_append(relate_from.via, r.id) as via
			from relations r
			inner join relate_from on relate_from."to" = r."from"
			where r."to" != $1
		) select * from relate_from order by id
	`

	var relations []RecRelation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func RetrieveRelation(ctx context.Context, id string) (*Relation, error) {
	query := `
		select
		where id = $1
	`

	var r Relation

	if err := pgxscan.Get(ctx, pg, &r, query, id); err != nil {
		return nil, fmt.Errorf("get failed: %w", err)
	}

	return &r, nil
}

func (r *Relation) Delete(ctx context.Context) error {
	query := `delete from relations where id = $1`

	if _, err := pg.Exec(ctx, query, r.ID); err != nil {
		return fmt.Errorf("pg.Exec failed: %w", err)
	}

	caches, err := listDerivativeCaches(ctx, pg, *r)
	if err != nil {
		return fmt.Errorf("listDerivativeCaches failed: %w", err)
	}

	for _, cache := range caches {
		fmt.Println("del", cache.ID)
		if err := (&cache).Delete(ctx); err != nil {
			return fmt.Errorf("cache.Delete failed: %w", err)
		}
	}

	return nil
}
