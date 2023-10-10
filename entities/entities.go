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

type CreateRequest struct {
	ID   string
	Type string

	Attrs map[string]any
}

type ListRequest struct {
	Type            *string
	PaginationToken *string
}

type ListResponse struct {
	Data            []Entity
	PaginationToken *string
}

type UpdateRequest struct {
	ID   string
	Type string

	Attrs map[string]any
}

func Get(ctx context.Context, id, typ string) (*Entity, error) {
	dbentity, err := db.Get(ctx, id, typ)
	if err != nil {
		return nil, err
	}

	return u.Ptr(mapFromDB(*dbentity)), nil
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

	return u.Ptr(mapFromDB(*entity)), nil
}

func Create(ctx context.Context, request CreateRequest) (*Entity, error) {
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

	return u.Ptr(mapFromDB(*dbe)), nil
}

func List(ctx context.Context, r ListRequest) (*ListResponse, error) {
	entities, err := db.List(ctx, db.Filter{
		Type:            r.Type,
		PaginationToken: r.PaginationToken,
	})
	if err != nil {
		return nil, err
	}

	res := ListResponse{Data: u.Map(entities, mapFromDB)}
	if len(res.Data) > 0 {
		res.PaginationToken = &res.Data[len(res.Data)-1].ID
	}

	return &res, err
}

func Delete(ctx context.Context, id, typ string) error {
	dbentity := &db.Entity{ID: id, Type: typ}
	if err := dbentity.Delete(ctx); err != nil {
		return err
	}
	return nil
}

func mapFromDB(e db.Entity) Entity {
	return Entity{
		ID:        e.ID,
		Type:      e.Type,
		Attrs:     e.Attrs,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
