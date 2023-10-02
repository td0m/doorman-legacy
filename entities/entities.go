package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/u"
)

type Entity struct {
	ID   string
	Type string // todo enum

	Attrs     map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *Entity) EntityID() string {
	return e.ID
}

func (e *Entity) EntityType() string {
	return e.Type
}

type UpdateRequest struct {
	ID   string
	Type string

	Attrs map[string]any
}

type EntityUpdate struct {
	Attrs map[string]any
}

func Get(ctx context.Context, id, typ string) (*Entity, error) {
	panic(3)
}

func Update(ctx context.Context, request UpdateRequest) (*Entity, error) {
	entity, err := db.Get(ctx, request.ID, request.Type)
	if err != nil {
		return nil, err
	}

	if request.Attrs != nil {
		entity.Attrs = request.Attrs
	}

	if err := entity.Update(ctx); err != nil {
		return nil, fmt.Errorf("Update failed: %w", err)
	}

	return u.Ptr(toDomain(*entity)), nil
}

func Create(ctx context.Context, request Entity) (*Entity, error) {
	if request.Type == "" {
		return nil, fmt.Errorf("type is required")
	}
	if request.ID == "" {
		request.ID = xid.New().String()
	}
	dbe := &db.Entity{
		ID:    request.ID,
		Type:  request.Type,
		Attrs: request.Attrs,
	}

	if err := dbe.Create(ctx); err != nil {
		return nil, err
	}

	request = toDomain(*dbe)
	return &request, nil
}

func Delete(ctx context.Context, id string) error {
	return nil
}

func toDomain(e db.Entity) Entity {
	return Entity{
		ID:        e.ID,
		Type:      e.Type,
		Attrs:     e.Attrs,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
