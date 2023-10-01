package relations

import (
	"context"
	"fmt"

	"github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
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

func Create(ctx context.Context, req CreateRequest) (*Relation, error) {
	fmt.Printf("%+v", req)
	dbrelation := &db.Relation{
		ID:    req.ID,
		From:  entityToDB(req.From),
		To:    entityToDB(req.To),
		Attrs: req.Attrs,
	}

	if err := dbrelation.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed to create: %w", err)
	}

	relation := toDomain(*dbrelation)

	return &relation, nil
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
