package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/rs/xid"
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
	From string
	Name string
	To   string
}

type Dependency struct {
	RelationID string `db:"relation_id"`
	DependsOn  string `db:"depends_on"`
}

func (d *Dependency) Create(ctx context.Context) error {
	query := `
		insert into dependencies(relation_id, depends_on)
		values($1, $2)
	`

	if _, err := pg.Exec(ctx, query, d.RelationID, d.DependsOn); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (r *Relation) Create(ctx context.Context) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	query := `
		insert into relations(id, "from", name, "to")
		values($1, $2, $3, $4)
	`

	if _, err := pg.Exec(ctx, query, r.ID, r.From, r.Name, r.To); err != nil {
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

		for _, to := range tos {
			if extractType(to.To) != "user" {
				continue
			}
			fmt.Println("got", extractType(to.To))
			for _, from := range froms {
				if extractType(from.To) != "collection" {
					continue
				}
				if from.From == r.From && to.To == r.To {
					fmt.Println("ignoring self")
					continue
				}
				id := xid.New().String()
				if _, err := pg.Exec(ctx, query, id, from.From, from.Name, to.To); err != nil {
					return fmt.Errorf("creating derivative relation failed: %w", err)
				}
				deps := []string{}
				deps = append(deps, from.Via...)
				deps = append(deps, to.Via...)
				for _, dependsOn := range deps {
					dep := &Dependency{RelationID: id, DependsOn: dependsOn}
					if err := dep.Create(ctx); err != nil {
						return fmt.Errorf("failed to create dep: (%+v): %w", dep, err)
					}
				}
			}
		}
	}

	return nil
}

func extractType(s string) string {
	return strings.SplitN(s, ":", 2)[0]
}

func RetrieveRelation(ctx context.Context, id string) (*Relation, error) {
	query := `
		select id, "from", name, "to"
		from relations
		where id = $1
	`

	var r Relation
	if err := pgxscan.Get(ctx, pg, &r, query, id); err != nil {
		return nil, err
	}

	return &r, nil
}

func listRecRelationsTo(ctx context.Context, tx pgxscan.Querier, id string) ([]RecRelation, error) {
	query := `
		with recursive relate_to as(
			select
				id,
				"from",
				"to",
				name,
				array_append('{}'::text[], id) as via
			from relations
			where "to" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				r.name,
				array_append(relate_to.via, r.id) as via
			from relations r
			inner join relate_to on relate_to."from" = r."to"
			where r."from" != $1
		) select id, "from", name, "to", via from relate_to order by id desc
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
				name,
				array_append('{}'::text[], id) as via
			from relations
			where "from" = $1

			union

			select
				r.id,
				r."from",
				r."to",
				r.name,
				array_append(relate_from.via, r.id) as via
			from relations r
			inner join relate_from on relate_from."to" = r."from"
			where r."to" != $1
		) select id, "from", name, "to", via from relate_from order by id
	`

	var relations []RecRelation
	err := pgxscan.Select(ctx, tx, &relations, query, id)
	if err != nil {
		return nil, fmt.Errorf("select failed: %w", err)
	}

	return relations, nil
}

func listDerivatives(ctx context.Context, tx pgxscan.Querier, r Relation) ([]RecRelation, []RecRelation, error) {
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
	// for _, to := range to {
	// 	if to.To == r.From {
	// 		return ErrCycle
	// 	}
	// }

	froms = append(froms, RecRelation{Relation: r})
	tos = append(tos, RecRelation{Relation: r})
	return froms, tos, nil
}
