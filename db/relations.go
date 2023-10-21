package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrCycle = errors.New("cycle detected")
var ErrRelationExists = errors.New("this relation already exists")

type Relation struct {
	From string
	Name string
	To   string

	Via []string
}

func Check(ctx context.Context, from, name, to string) ([]Relation, error) {
	query := `
		select "from", "to", name
		from relations
		where
			"from" = $1 and
			name   = $2 and
			"to"   = $3 `

	var relations []Relation
	if err := pgxscan.Select(ctx, pg, &relations, query, from, name, to); err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func (r *Relation) Delete(ctx context.Context) error {
	if r.Via == nil {
		r.Via = []string{}
	}
	query := `
		delete from relations
		where
			"from" = $1 and
			"name" = $2 and
			"to"   = $3 and
			via    = $4
	`

	if _, err := pg.Exec(ctx, query, r.From, r.Name, r.To, r.Via); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	froms, tos, err := listDerivatives(ctx, pg, *r)
	if err != nil {
		return fmt.Errorf("listDerivatives failed: %w", err)
	}

	// TODO: remove all
	rels := flattenedCollections(*r, froms, tos)
	for _, r := range rels {
		fmt.Printf("del: %+v\n", r)
		if _, err := pg.Exec(ctx, query, r.From, r.Name, r.To, r.Via); err != nil {
			return fmt.Errorf("deleting derivative relation failed: %w", err)
		}
	}

	return nil
}

func (r *Relation) Create(ctx context.Context) error {
	query := `
		insert into relations("from", name, "to", via)
		values($1, $2, $3, $4)
	`

	if _, err := pg.Exec(ctx, query, r.From, r.Name, r.To, []string{}); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "relations_pkey" {
				return ErrRelationExists
			}
		}
		return fmt.Errorf("exec failed: %w", err)
	}

	// TODO: always ensure no cycles

	fromType := strings.SplitN(r.From, ":", 2)[0]
	toType := strings.SplitN(r.To, ":", 2)[0]
	if fromType == "collection" || toType == "collection" {
		froms, tos, err := listDerivatives(ctx, pg, *r)
		if err != nil {
			return fmt.Errorf("listDerivatives failed: %w", err)
		}

		changes := flattenedCollections(*r, froms, tos)
		for _, r := range changes {
			if _, err := pg.Exec(ctx, query, r.From, r.Name, r.To, r.Via); err != nil {
				return fmt.Errorf("creating derivative relation failed: %w", err)
			}
		}
	}

	return nil
}

func extractType(s string) string {
	return strings.SplitN(s, ":", 2)[0]
}

func ListForward(ctx context.Context, from, name string) ([]Relation, error) {
	query := `
		select "from", name, "to"
		from relations
		where
			"from" = $1 and
			name   = $2
	`

	var rs []Relation
	if err := pgxscan.Select(ctx, pg, &rs, query, from, name); err != nil {
		return nil, err
	}

	return rs, nil
}

func listRecRelationsTo(ctx context.Context, tx pgxscan.Querier, id string) ([]Relation, error) {
	query := `
		with recursive relate_to as(
			select
				"from",
				"to",
				name,
				array_append(array_append('{}'::text[], "from"), name) as via
			from relations
			where "to" = $1

			union

			select
				r."from",
				r."to",
				r.name,
				array_prepend(r.from, array_prepend(r.name, relate_to.via)) as via
			from relations r
			inner join relate_to on relate_to."from" = r."to"
			where r."from" != $1
		) select "from", name, "to", via from relate_to
	`

	var relations []Relation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func listRecRelationsFrom(ctx context.Context, tx pgxscan.Querier, id string) ([]Relation, error) {
	query := `
		with recursive relate_from as(
			select
				"from",
				"to",
				name,
				array_append(array_append('{}'::text[], name), "to") as via
			from relations
			where "from" = $1

			union

			select
				r."from",
				r."to",
				r.name,
				array_append(array_append(relate_from.via, r.name), r.to) as via
			from relations r
			inner join relate_from on relate_from."to" = r."from"
			where r."to" != $1
		) select "from", name, "to", via from relate_from
	`

	var relations []Relation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func listDerivatives(ctx context.Context, tx pgxscan.Querier, r Relation) ([]Relation, []Relation, error) {
	froms, err := listRecRelationsTo(ctx, tx, r.From)
	if err != nil {
		return nil, nil, fmt.Errorf("listRecRelationsTo failed: %w", err)
	}
	tos, err := listRecRelationsFrom(ctx, tx, r.To)
	if err != nil {
		return nil, nil, fmt.Errorf("listRecRelationsFrom failed: %w", err)
	}

	// Because of the nature of cycles, this will always match.
	// No need for the second statement.
	for _, from := range froms {
		if from.From == r.To {
			return nil, nil, ErrCycle
		}
	}

	froms = append(froms, r)
	tos = append(tos, r)
	return froms, tos, nil
}

func flattenedCollections(r Relation, froms, tos []Relation) []Relation {
	out := []Relation{}

	for _, to := range tos {
		if extractType(to.To) != "user" {
			continue
		}
		for _, from := range froms {
			if extractType(from.To) != "collection" {
				continue
			}
			// NOT SURE IF I NEED THIS
			if from.From == r.From && to.To == r.To {
				fmt.Println("ignoring self")
				continue
			}

			deps := []string{}
			deps = append(deps, from.Via...)
			deps = append(deps, r.From, r.Name, r.To)
			deps = append(deps, to.Via...)
			fmt.Printf("deps %+v\n", deps)

			// TODO: figure out how to not duplicate?
			out = append(out, Relation{
				From: from.From,
				To:   to.To,
				Name: from.Name,
				Via:  deps,
			})
		}
	}
	return out
}
