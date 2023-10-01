package relations

import (
	"context"
	"fmt"

	"github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
)

type Relation struct {
	ID    string
	From  string
	To    string
	Attrs map[string]any
}

type CreateRequest struct {
	ID    string
	From  string
	To    string
	Attrs map[string]any
}

type ListRequest struct {
	From string
	To   string
}

func Entity(typ, id string) string {
	return typ + ":" + id
}

func Create(ctx context.Context, r CreateRequest) (*Relation, error) {
	dbrelation := &db.Relation{
		ID:    r.ID,
		From:  r.From,
		To:    r.To,
		Attrs: r.Attrs,
	}

	if err := dbrelation.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed to create: %w", err)
	}

	relation := toDomain(*dbrelation)

	return &relation, nil
}

func List(ctx context.Context, r ListRequest) ([]Relation, error) {
	if r.From == "" || r.To == "" {
		return nil, fmt.Errorf("to and from must be provided")
	}

	dbrelations, err := db.ListRelations(ctx, db.RelationFilter{
		From: &r.From,
		To: &r.To,
	})
	if err != nil {
		return nil, err
	}

	return u.Map(dbrelations, toDomain), nil
}

func toDomain(r db.Relation) Relation {
	return Relation{
		ID:    r.ID,
		From:  r.From,
		To:    r.To,
		Attrs: r.Attrs,
	}
}
