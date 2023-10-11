package entities

import (
	"context"

	"github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/u"
)

type CreateTypeRequest struct {
	ID    string
	Attrs map[string]any
}

type UpdateTypeRequest struct {
	ID    string
	Attrs map[string]any
}

type Type struct {
	ID    string
	Attrs map[string]any
}

type ListTypesResponse struct {
	Data []Type
}

func CreateType(ctx context.Context, r CreateTypeRequest) (*Type, error) {
	dbtype := &db.Type{ID: r.ID, Attrs: r.Attrs}
	if err := dbtype.Create(ctx); err != nil {
		return nil, err
	}

	return u.Ptr(mapTypeFromDB(*dbtype)), nil
}

func UpdateType(ctx context.Context, r UpdateTypeRequest) (*Type, error) {
	dbtype, err := db.RetrieveType(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	if r.Attrs != nil {
		dbtype.Attrs = r.Attrs
	}

	if err := dbtype.Update(ctx); err != nil {
		return nil, err
	}

	return u.Ptr(mapTypeFromDB(*dbtype)), nil
}

func ListTypes(ctx context.Context) (*ListTypesResponse, error) {
	dbtypes, err := db.ListTypes(ctx)
	if err != nil {
		return nil, err
	}

	return &ListTypesResponse{
		Data: u.Map(dbtypes, mapTypeFromDB),
	}, nil
}

func mapTypeFromDB(t db.Type) Type {
	return Type{ID: t.ID, Attrs: t.Attrs}
}
