package service

import (
	"context"
	"fmt"

	"github.com/td0m/poc-doorman/db"
)

type Entities struct{}

type Entity struct {
	ID    string
	Attrs map[string]any
}

type CreateEntity struct {
	ID    string
	Attrs map[string]any
}

func (es *Entities) Create(ctx context.Context, request CreateEntity) (*Entity, error) {
	e := &db.Entity{
		ID:    request.ID,
		Attrs: request.Attrs,
	}
	if err := e.Create(ctx); err != nil {
		return nil, fmt.Errorf("db.Create failed: %w", err)
	}

	res := mapEntityFromDB(*e)
	return &res, nil
}

func mapEntityFromDB(e db.Entity) Entity {
	return Entity{
		ID:    e.ID,
		Attrs: e.Attrs,
	}
}
