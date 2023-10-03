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

// anything not here can only be a leaf node
var validRelations = map[string][]string{
	"collection": {"collection", "role", "*"},
	"role":       {"permission"},
	"user":       {"collection", "role", "*"},
	"permission": {},
}

type Relation struct {
	ID   string
	From Entity
	To   Entity
	Name *string
}

type Entity struct {
	ID   string
	Type string
}

type CreateRequest struct {
	ID   string
	From Entity
	To   Entity
	Name *string
}

type ListRequest struct {
	From *Entity
	To   *Entity
}

type UpdateRequest struct {
	Name *string
}

func Create(ctx context.Context, req CreateRequest) (*Relation, error) {
	if req.Name != nil && !canHaveName(req.From.Type, req.To.Type) {
		return nil, errs.New(http.StatusBadRequest, "cannot set name for this relation: invalid types")
	}
	dbrelation := &db.Relation{
		ID:   req.ID,
		From: entityToDB(req.From),
		To:   entityToDB(req.To),
		Name: req.Name,
	}

	canConnectTo, ok := validRelations[req.From.Type]
	if !ok {
		return nil, errs.New(http.StatusBadRequest, "invalid relation")
	}

	to := req.To.Type
	if _, toTypeIsStrict := validRelations[req.To.Type]; !toTypeIsStrict {
		to = "*"
	}

	if canConnect := slices.Contains(canConnectTo, to); !canConnect {
		return nil, errs.New(http.StatusBadRequest, "cannot connect to this type")
	}

	if err := dbrelation.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed to create: %w", err)
	}

	dbrelation.Cache = true
	if err := dbrelation.Create(ctx); err != nil {
		return nil, fmt.Errorf("db failed to create cache: %w", err)
	}

	relation := toDomain(*dbrelation)

	if err := rebuildCache(ctx, relation, req.From, req.To); err != nil {
		return nil, fmt.Errorf("failed to rebuild cache: %w", err)
	}

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

	dbrelations, err := db.ListRelations(ctx, filter, true) // todo: use cache prop, true by default
	if err != nil {
		return nil, err
	}

	return u.Map(dbrelations, toDomain), nil
}

func Delete(ctx context.Context, id string) error {
	if err := u.Ptr(db.Relation{ID: id}).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context, id string) (*Relation, error) {
	dbrelation, err := db.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.Ptr(toDomain(*dbrelation)), nil
}

func toDomain(r db.Relation) Relation {
	return Relation{
		ID:   r.ID,
		From: entityRefToDomain(r.From),
		To:   entityRefToDomain(r.To),
		Name: r.Name,
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

// This is where the magic happens!
func rebuildCache(ctx context.Context, relation Relation, from, to Entity) error {
	leftRelations, err := db.ListLinkedToWithDependencies(ctx, entityToDB(from))
	if err != nil {
		return fmt.Errorf("failed to list relations to the left: %w", err)
	}
	leftRelations = append(leftRelations, db.RelationWithDeps{
		Relation: db.Relation{From: entityToDB(from)},
	})

	rightRelations, err := db.ListLinkedFromWithDependencies(ctx, entityToDB(to))
	if err != nil {
		return fmt.Errorf("failed to list relations to the right: %w", err)
	}
	rightRelations = append(rightRelations, db.RelationWithDeps{
		Relation: db.Relation{To: entityToDB(to)},
	})

	// Derrive name if set by any of the relations
	var name *string
	for _, rel := range leftRelations {
		if rel.Name != nil {
			name = rel.Name
			break
		}
	}
	if name == nil {
		for _, rel := range rightRelations {
			if rel.Name != nil {
				name = rel.Name
				break
			}
		}
	}

	for i, l := range leftRelations {
		for j, r := range rightRelations {
			// Skip the new relation as it was created outside of this function
			if i == len(leftRelations)-1 && j == len(rightRelations)-1 {
				continue
			}

			cache := &db.Relation{
				From:  l.From,
				To:    r.To,
				Cache: true,
				Name:  name,
			}
			if err := cache.Create(ctx); err != nil {
				return fmt.Errorf("failed to create rel: %w", err)
			}

			deps := []string{}
			if l.ID != "" {
				if len(l.DependencyIDs) == 0 {
					deps = append(deps, l.ID) // Direct! Add as dependency
				} else {
					deps = append(deps, l.DependencyIDs...) // Cache, so copy dependencies!
				}
			}

			deps = append(deps, relation.ID)

			if r.ID != "" {
				if len(r.DependencyIDs) == 0 {
					deps = append(deps, r.ID) // Direct! Add as dependency
				} else {
					deps = append(deps, r.DependencyIDs...) // Cache, so copy dependencies!
				}
			}

			for _, dep := range deps {
				dbdep := &db.Dependency{
					CacheID:    cache.ID,
					RelationID: dep,
				}
				if err := dbdep.Create(ctx); err != nil {
					return fmt.Errorf("failed to create dep (%+v): %w", dbdep, err)
				}
			}
		}
	}
	return nil
}

func canHaveName(from, to string) bool {
	return (from == "collection" || from == "user") && to != "collection"
}
