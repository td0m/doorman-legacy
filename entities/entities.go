package entities

import (
	"context"
	"fmt"

	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/entities/db"
)

type Entity struct {
	ID   string
	Type string // todo enum

	Attrs map[string]any
}

func (e *Entity) EntityID() string {
	return e.ID
}

func (e *Entity) EntityType() string {
	return e.Type
}

type EntityUpdate struct {
	Attrs map[string]any
}

func Get(ctx context.Context, eType, id string) (*Entity, error) {
	panic(3)
}

func Update(ctx context.Context, eType, id string, changes EntityUpdate) (*Entity, error) {
	panic(3)
}

func Create(ctx context.Context, e Entity) (*Entity, error) {
	if e.Type == "" {
		return nil, fmt.Errorf("type is required")
	}
	if e.ID == "" {
		e.ID = xid.New().String()
	}
	dbe := &db.Entity{
		ID:   e.ID,
		Type: e.Type,
	}

	if err := dbe.Create(ctx); err != nil {
		return nil, err
	}

	e = toDomain(*dbe)
	return &e, nil
}

func Delete(ctx context.Context, id string) error {
	return nil
}

func toDomain(e db.Entity) Entity {
	return Entity{
		ID:   e.ID,
		Type: e.Type,
	}
}
