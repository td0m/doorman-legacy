package relations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/td0m/poc-doorman/errs"
	"github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
	"golang.org/x/exp/slices"
)

type Relation struct {
	ID    string
	From  Entity
	To    Entity
	Attrs map[string]any
}

type Entity struct {
	ID   string
	Type string
}

type CreateRequest struct {
	ID    string
	From  Entity
	To    Entity
	Attrs map[string]any
}

type ListRequest struct {
	From *Entity
	To   *Entity
}

// func Entity(typ, id string) string {
// 	return typ + ":" + id
// }

var validRelations = map[string][]string{
	"collection": {"role", "resource"},
	"permission": {},
	"resource":   {},
	"role":       {"permission"},
	"user":       {"collection", "role", "resource"},
}

func Create(ctx context.Context, req CreateRequest) (*Relation, error) {
	dbrelation := &db.Relation{
		ID:    req.ID,
		From:  entityToDB(req.From),
		To:    entityToDB(req.To),
		Attrs: req.Attrs,
	}

	canConnectTo, ok := validRelations[req.From.Type]
	if !ok {
		return nil, errs.New(http.StatusBadRequest, "invalid relation")
	}

	if canConnect := slices.Contains(canConnectTo, req.To.Type); !canConnect {
		return nil, errs.New(http.StatusBadRequest, "cannot connect to this type")
	}

	if err := rebuildIndirects(ctx, req.From, req.To); err != nil {
		return nil, fmt.Errorf("failed to rebuild indirects: %w", err)
	}

	if err := dbrelation.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed to create: %w", err)
	}

	relation := toDomain(*dbrelation)

	return &relation, nil
}

func rebuildIndirects(ctx context.Context, from, to Entity) error {
	fmt.Println("rebuild:", from.Type, "=>", to.Type)

	leftRelations, err := db.ListRelations(ctx, db.RelationFilter{
		ToID:   &from.ID,
		ToType: &from.Type,
	})
	if err != nil {
		return fmt.Errorf("failed to list relations to the left: %w", err)
	}
	left := make([]Entity, len(leftRelations))
	for i, r := range leftRelations {
		left[i] = entityRefToDomain(r.From)
	}
	left = append(left, from)

	rightRelations, err := db.ListRelations(ctx, db.RelationFilter{
		FromID:   &to.ID,
		FromType: &to.Type,
	})
	if err != nil {
		return fmt.Errorf("failed to list relations to the right: %w", err)
	}
	right := make([]Entity, len(rightRelations))
	for i, r := range rightRelations {
		right[i] = entityRefToDomain(r.To)
	}
	right = append(right, to)

	for _, l := range left {
		for _, r := range right {
			if l == from && r == to {
				continue
			}
			rel := &db.Relation{
				From:     entityToDB(l),
				To:       entityToDB(r),
				Indirect: true,
			}
			if err := rel.Create(ctx); err != nil {
				return fmt.Errorf("failed to create rel: %w", err)
			}
		}
	}
	return nil
}

func List(ctx context.Context, r ListRequest) ([]Relation, error) {
	if r.From == nil && r.To == nil {
		return nil, fmt.Errorf("to or from must be provided")
	}

	filter := db.RelationFilter{}
	if r.From != nil {
		filter.FromID = &r.From.ID
		filter.FromType = &r.From.Type
	}
	if r.To != nil {
		filter.ToID = &r.To.ID
		filter.ToType = &r.To.Type
	}

	dbrelations, err := db.ListRelations(ctx, filter)
	if err != nil {
		return nil, err
	}

	return u.Map(dbrelations, toDomain), nil
}

func toDomain(r db.Relation) Relation {
	return Relation{
		ID:    r.ID,
		From:  entityRefToDomain(r.From),
		To:    entityRefToDomain(r.To),
		Attrs: r.Attrs,
	}
}

func entityToDB(r Entity) db.EntityRef {
	return db.EntityRef{
		ID:   r.ID,
		Type: r.Type,
	}
}

func entityRefToDomain(r db.EntityRef) Entity {
	return Entity{
		ID:   r.ID,
		Type: r.Type,
	}
}
