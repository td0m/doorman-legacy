package service

import (
	"context"
	"fmt"

	"github.com/td0m/poc-doorman/db"
)

type Relation struct {
	ID   string
	Name *string
	From string
	To   string
}
type Relations struct {
}

type RelationsListRequest struct {
	NoCache         bool
	From            *string
	FromType        *string
	To              *string
	ToType          *string
	Name            *string
	PaginationToken *string
}

type RelationList struct {
	Items           []Relation
	PaginationToken string
}

type RelationsCreate struct {
	From string
	To   string
	Name *string
}

func (*Relations) Create(ctx context.Context, request RelationsCreate) (*Relation, error) {
	r := &db.Relation{
		From: request.From,
		To:   request.To,
		Name: request.Name,
	}
	if err := r.Create(ctx); err != nil {
		return nil, fmt.Errorf("create failed: %w", err)
	}

	res := mapRelationFromDB(*r)
	return &res, nil
}

func mapRelationFromDB(r db.Relation) Relation {
	return Relation{
		ID:   r.ID,
		From: r.From,
		To:   r.To,
	}
}

func mapRelationFromDBCache(r db.Cache) Relation {
	return Relation{
		ID:   r.ID,
		From: r.From,
		To:   r.To,
	}
}
func (*Relations) Delete(ctx context.Context, id string) error {
	r, err := db.RetrieveRelation(ctx, id)
	if err != nil {
		return fmt.Errorf("db.RetrieveRelation failed: %w", err)
	}

	if err := r.Delete(ctx); err != nil {
		return fmt.Errorf("reaction.Delete failed: %w", err)
	}
	return nil
}

// GET /relations?no_cache=true&from=user:alice&to_type=collection
func (*Relations) List(ctx context.Context, request RelationsListRequest) (RelationList, error) {
	table := "cache"
	if request.NoCache {
		table = "relations"
	}

	f := db.RelationFilter{
		AfterID:  request.PaginationToken,
		From:     request.From,
		FromType: request.FromType,
		Name:     request.Name,
		To:       request.To,
		ToType:   request.ToType,
	}

	relations, err := db.ListRelationsOrCache(ctx, table, f)
	if err != nil {
		return RelationList{}, fmt.Errorf("db failed: %w", err)
	}

	items := make([]Relation, len(relations))
	for i := range relations {
		items[i] = mapRelationFromDBCache(relations[i])
	}

	return RelationList{
		Items: items,
	}, nil
}
