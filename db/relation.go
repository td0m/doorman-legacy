package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
	"golang.org/x/exp/slog"
)

var ErrCycle = errors.New("cycle detected")

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

	// TODO: validate types
	// fromType, toType := extractType(r.From), extractType(r.To)

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
			slog.Error("tx.Rollback failed in relation.Create: ", err)
		}
	}()

	// Concurrent adding and removing relations can lead to dangling cache without this
	if _, err := tx.Exec(ctx, `lock table relations in access exclusive mode`); err != nil {
		return fmt.Errorf("locking table relations failed: %w", err)
	}
	query := `
		insert into relations(id, "from", "to", name)
		values($1, $2, $3, $4)
	`

	_, err = tx.Exec(ctx, query, r.ID, r.From, r.To, r.Name)
	if err != nil {
		return fmt.Errorf("tx.Exec failed to insert relation: %w", err)
	}

	// Computes relations
	{
		// TODO: can we use these two to prevent cycles?
		from, err := listRecRelationsTo(ctx, tx, r.From)
		if err != nil {
			return fmt.Errorf("listRecRelationsTo failed: %w", err)
		}
		to, err := listRecRelationsFrom(ctx, tx, r.To)
		if err != nil {
			return fmt.Errorf("listRecRelationsFrom failed: %w", err)
		}

		// Because of the nature of cycles, this will always match.
		// No need for the second statement.
		for _, from := range from {
			if from.From == r.To {
				return ErrCycle
			}
		}
		// for _, to := range to {
		// 	if to.To == r.From {
		// 		return ErrCycle
		// 	}
		// }

		from = append(from, RecRelation{Relation: *r})
		to = append(to, RecRelation{Relation: *r})

		for _, from := range from {
			for _, to := range to {
				cache := &Cache{
					Via:  append(append(from.Via, r.ID), to.Via...),
					From: from.From,
					To:   to.To,
					Name: nil,
				}

				if err := cache.Create(ctx, tx); err != nil {
					return fmt.Errorf("cache.Create failed: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("tx.Commit failed: %w", err)
	}

	return nil
}

func listRecRelationsTo(ctx context.Context, tx pgx.Tx, id string) ([]RecRelation, error) {
	query := `
		with recursive relate_to as(
			select
				id,
				"from",
				"to",
				'{}'::text[] as via
			from relations
			where "to" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				array_append(relate_to.via, relate_to.id) as via
			from relations r
			inner join relate_to on relate_to."from" = r."to"
			where r."from" != $1
		) select * from relate_to
	`

	var relations []RecRelation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func listRecRelationsFrom(ctx context.Context, tx pgx.Tx, id string) ([]RecRelation, error) {
	query := `
		with recursive relate_from as(
			select
				id,
				"from",
				"to",
				'{}'::text[] as via
			from relations
			where "from" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				array_append(relate_from.via, relate_from.id) as via
			from relations r
			inner join relate_from on relate_from."to" = r."from"
			where r."to" != $1
		) select * from relate_from
	`

	var relations []RecRelation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

// TODO: remove.
