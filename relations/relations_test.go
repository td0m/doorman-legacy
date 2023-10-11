package relations

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	entitiesdb "github.com/td0m/poc-doorman/entities/db"
	"github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

const n int = 1_000_000

func setupSampleData() {
	ctx := context.Background()
	// TODO: just use generate_series
	fmt.Println("Creating users and resources...", time.Now())
	for t := 0; t < n/10; t++ {
		var g errgroup.Group
		for i := 0; i < 0; i++ {
			i := t*10 + i
			g.Go(func() error {
				e := &entitiesdb.Entity{
					Type: "resource",
					ID:   strconv.Itoa(i),
				}
				return e.Create(ctx)
			})
			g.Go(func() error {
				user := &entitiesdb.Entity{
					Type: "user",
					ID:   strconv.Itoa(i),
				}
				return user.Create(ctx)
			})
		}
		u.Check(g.Wait())
	}

	// fmt.Println("Creating relationships...")
	// for t := 0; t < n; t++ {
	// 	var g errgroup.Group
	// 	for i := 0; i < 10; i++ {
	// 		g.Go(func() error {
	// 			in := CreateRequest{
	// 				From: Entity{Type: "user", ID: strconv.Itoa(rand.Intn(n - 1))},
	// 				To:   Entity{Type: "resource", ID: strconv.Itoa(rand.Intn(n - 1))},
	// 				Name: u.Ptr("owner"),
	// 			}
	// 			_, err := Create(ctx, in)
	// 			return err
	// 		})
	// 	}
	// 	u.Check(g.Wait())
	// }

	for i := 0; i < 100000; i++ {
		start := time.Now()
		params := []any{}
		query := strings.Builder{}

		query.WriteString(`
	  insert into cache(_id, from_id, from_type, to_id, to_type, name)
		values
	`)
		m := 1000
		for i := 0; i < m; i++ {
			row := []any{
				xid.New().String(),
				strconv.Itoa(rand.Intn(n - 1)),
				"user",
				strconv.Itoa(rand.Intn(n - 1)),
				"resource",
				"owner",
			}
			args := []string{}
			for i := range row {
				args = append(args, "$"+strconv.Itoa(i+len(params)+1))
			}
			params = append(
				params,
				row...,
			)
			if i > 0 {
				query.WriteString(",")
			}
			query.WriteString("(" + strings.Join(args, ",") + ")")
		}

		_, err := db.Conn.Exec(ctx, query.String(), params...)
		u.Check(err)
		fmt.Println(time.Since(start).Nanoseconds() / int64(m))
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)
	defer pgDoorman.Close()

	entitiesdb.Conn = pgDoorman
	db.Conn = pgDoorman

	_ = (&entitiesdb.Type{ID: "resource"}).Create(ctx)
	_ = (&entitiesdb.Type{ID: "post"}).Create(ctx)

	// setupSampleData()

	m.Run()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	t.Run("Success on valid entities", func(t *testing.T) {
		user := &entitiesdb.Entity{
			Type: "user",
			ID:   xid.New().String(),
		}
		require.NoError(t, user.Create(ctx))

		resource := &entitiesdb.Entity{
			Type: "resource",
			ID:   xid.New().String(),
		}
		require.NoError(t, resource.Create(ctx))

		req := CreateRequest{
			ID:   xid.New().String(),
			From: Entity{ID: user.ID, Type: user.Type},
			To:   Entity{ID: resource.ID, Type: resource.Type},
		}
		relation, err := Create(ctx, req)
		require.NoError(t, err)
		require.Equal(t, req.ID, relation.ID)
		require.Equal(t, req.From.ID, relation.From.ID)
		require.Equal(t, req.From.Type, relation.From.Type)
		require.Equal(t, req.To.ID, relation.To.ID)
		require.Equal(t, req.To.Type, relation.To.Type)
		require.Equal(t, req.Name, relation.Name)
	})

	t.Run("Failure on connection to self", func(t *testing.T) {
		coll1 := &entitiesdb.Entity{
			Type: "collection",
		}
		require.NoError(t, coll1.Create(ctx))

		req := CreateRequest{
			ID:   xid.New().String(),
			From: Entity{ID: coll1.ID, Type: coll1.Type},
			To:   Entity{ID: coll1.ID, Type: coll1.Type},
		}
		_, err := Create(ctx, req)
		assert.Error(t, err)
	})

	t.Run("Failure on missing \"from\" entity", func(t *testing.T) {
		user := &entitiesdb.Entity{
			Type: "user",
			ID:   xid.New().String(),
		}
		require.NoError(t, user.Create(ctx))

		req := CreateRequest{
			ID:   xid.New().String(),
			From: Entity{ID: user.ID, Type: user.Type},
		}
		_, err := Create(ctx, req)
		require.Error(t, err)
	})

	t.Run("Validates entity type", func(t *testing.T) {
		tests := []struct {
			FromType string
			ToType   string
			Success  bool
		}{
			{"collection", "collection", true},
			{"collection", "permission", false},
			{"collection", "resource", true},
			{"collection", "role", true},
			{"collection", "user", false},
			{"permission", "collection", false},
			{"permission", "permission", false},
			{"permission", "resource", false},
			{"permission", "role", false},
			{"permission", "user", false},
			{"resource", "collection", false},
			{"resource", "permission", false},
			{"resource", "resource", false},
			{"resource", "role", false},
			{"resource", "user", false},
			{"role", "collection", false},
			{"role", "permission", true},
			{"role", "resource", false},
			{"role", "role", false},
			{"role", "user", false},
			{"user", "collection", true},
			{"user", "permission", false},
			{"user", "resource", true},
			{"user", "role", true},
			{"user", "user", false},

			{"user", "post", true},
			{"collection", "post", true},
			{"role", "post", false},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("Relating %s to %s results in sucess=%v", tt.FromType, tt.ToType, tt.Success), func(t *testing.T) {
				from := &entitiesdb.Entity{
					Type: tt.FromType,
					ID:   "id" + xid.New().String(),
				}
				require.NoError(t, from.Create(ctx))

				to := &entitiesdb.Entity{
					Type: tt.ToType,
					ID:   "id" + xid.New().String(),
				}
				require.NoError(t, to.Create(ctx))

				req := CreateRequest{
					ID:   xid.New().String(),
					From: Entity{ID: from.ID, Type: from.Type},
					To:   Entity{ID: to.ID, Type: to.Type},
				}
				_, err := Create(ctx, req)
				if tt.Success {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			})
		}
	})

	t.Run("Only allows names in certain relations", func(t *testing.T) {
		tests := []struct {
			from    string
			to      string
			success bool
		}{
			{"user", "collection", false},
			{"role", "permission", false},
			{"user", "permission", false},

			{"collection", "role", true},
			{"user", "role", true},
			{"user", "role", true},

			{"user", "post", true},
			{"collection", "post", true},
			{"collection", "resource", true},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s %s %v", tt.from, tt.to, tt.success), func(t *testing.T) {
				from := &entitiesdb.Entity{ID: "id" + xid.New().String(), Type: tt.from}
				require.NoError(t, from.Create(ctx))

				to := &entitiesdb.Entity{ID: "id" + xid.New().String(), Type: tt.to}
				require.NoError(t, to.Create(ctx))

				in := CreateRequest{
					From: Entity{ID: from.ID, Type: from.Type},
					To:   Entity{ID: to.ID, Type: to.Type},
					Name: u.Ptr("foo"),
				}
				_, err := Create(ctx, in)
				if tt.success {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("Sucess computing indirect relations", func(t *testing.T) {
		u1 := Entity{ID: "u1", Type: "user"}
		c1 := Entity{ID: "c1", Type: "collection"}
		r1 := Entity{ID: "r1", Type: "role"}
		p1 := Entity{ID: "p1", Type: "permission"}

		u1c1 := Relation{From: u1, To: c1}
		c1r1 := Relation{From: c1, To: r1}
		r1p1 := Relation{From: r1, To: p1}

		type relationWithDeps struct {
			Relation
			Deps []Relation
		}

		tests := []struct {
			Entities            []Entity
			Relations           []Relation
			AdditionalRelations []relationWithDeps
		}{
			{
				Entities:  []Entity{u1, c1, r1},
				Relations: []Relation{u1c1, c1r1},
				AdditionalRelations: []relationWithDeps{
					{Relation: Relation{From: u1, To: r1}, Deps: []Relation{u1c1, c1r1}},
				},
			},
			{
				Entities:  []Entity{u1, c1, r1, p1},
				Relations: []Relation{u1c1, c1r1, r1p1},
				AdditionalRelations: []relationWithDeps{
					{Relation: Relation{From: u1, To: r1}, Deps: []Relation{u1c1, c1r1}},
					{Relation: Relation{From: c1, To: p1}, Deps: []Relation{c1r1, r1p1}},
					{Relation: Relation{From: u1, To: p1}, Deps: []Relation{u1c1, c1r1, r1p1}},
				},
			},
		}
		for i, tt := range tests {
			t.Run(fmt.Sprintf("i %+v", i), func(t *testing.T) {
				permutations := permutations(tt.Relations)
				for _, relations := range permutations {
					// Insert them from left to right
					t.Run(fmt.Sprintf("Relations %+v", relations), func(t *testing.T) {
						rnd := xid.New().String()
						id := func(eid string) string {
							return "id" + eid + rnd
						}
						relationId := func(r Relation) string {
							return id(r.From.ID) + " => " + id(r.To.ID)
						}
						for i := range tt.Entities {
							e := tt.Entities[i]
							en := &entitiesdb.Entity{
								Type: e.Type,
								ID:   id(e.ID),
							}
							require.NoError(t, en.Create(ctx))
						}

						for _, pair := range relations {
							req := CreateRequest{
								ID:   relationId(pair),
								From: Entity{ID: id(pair.From.ID), Type: pair.From.Type},
								To:   Entity{ID: id(pair.To.ID), Type: pair.To.Type},
							}
							_, err := Create(ctx, req)
							require.NoError(t, err)

							res, err := List(ctx, ListRequest{
								From: &req.From,
								To:   &req.To,
							})
							require.NoError(t, err)
							require.Equal(t, 1, len(res.Data))
							// todo: check exists
						}

						tOnly := func(t Relation) string {
							return t.From.ID + "->" + t.To.ID
						}

						for _, rel := range tt.AdditionalRelations {
							req := ListRequest{
								From: &Entity{ID: id(rel.From.ID), Type: rel.From.Type},
								To:   &Entity{ID: id(rel.To.ID), Type: rel.To.Type},
							}
							all, err := List(ctx, req)
							require.NoError(t, err)
							require.Equal(t, 1, len(all.Data), "relation: %s, relations: %+v", rel.From.ID+" => "+rel.To.ID, u.Map(relations, tOnly))

							deps, err := db.ListDependencies(ctx, all.Data[0].ID)
							require.NoError(t, err)

							expectedDeps := u.Map(rel.Deps, relationId)

							slices.Sort(expectedDeps)
							slices.Sort(deps)

							require.Equal(t, expectedDeps, deps)
						}
					})
				}
			})
		}
	})
}

func TestListPagination(t *testing.T) {
	ctx := context.Background()

	createEntity := func(e Entity) error {
		return u.Ptr(entitiesdb.Entity{ID: e.ID, Type: e.Type}).Create(ctx)
	}

	rnd := xid.New().String()

	u1 := Entity{ID: "u1" + rnd, Type: "user"}
	require.NoError(t, createEntity(u1))

	for i := 0; i < 1_200; i++ {
		c := Entity{ID: "c" + strconv.Itoa(i) + rnd, Type: "collection"}
		require.NoError(t, createEntity(c))

		rel := db.Relation{
			From:  db.EntityRef{ID: u1.ID, Type: u1.Type},
			To:    db.EntityRef{ID: c.ID, Type: c.Type},
			Cache: true,
		}
		require.NoError(t, rel.Create(ctx))
	}

	t.Run("Only returns 1000 first items without pagination token", func(t *testing.T) {
		res, err := List(ctx, ListRequest{
			From: &u1,
		})
		require.NoError(t, err)
		assert.Equal(t, 1_000, len(res.Data))
		assert.Equal(t, "c999"+rnd, res.Data[999].To.ID)
		lastToken := res.PaginationToken
		t.Run("Second request returns remaining 200 items", func(t *testing.T) {
			res, err := List(ctx, ListRequest{
				From:            &u1,
				PaginationToken: lastToken,
			})
			require.NoError(t, err)
			assert.Equal(t, 200, len(res.Data))
			assert.Equal(t, "c1199"+rnd, res.Data[199].To.ID)
		})
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()

	rnd := xid.New().String()
	u1 := Entity{ID: "u1" + rnd, Type: "user"}
	c1 := Entity{ID: "c1" + rnd, Type: "collection"}

	createEntity := func(e Entity) error {
		return u.Ptr(entitiesdb.Entity{ID: e.ID, Type: e.Type}).Create(ctx)
	}

	require.NoError(t, createEntity(u1))
	require.NoError(t, createEntity(c1))

	rel := db.Relation{
		From:  db.EntityRef{ID: u1.ID, Type: u1.Type},
		To:    db.EntityRef{ID: c1.ID, Type: c1.Type},
		Name:  u.Ptr("foo"),
		Cache: true,
	}

	// Expect none before creation
	t.Run("Expect none before creation", func(t *testing.T) {
		rels, err := List(ctx, ListRequest{
			From: &u1,
			To:   &c1,
			Name: rel.Name,
		})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(rels.Data))
	})

	require.NoError(t, rel.Create(ctx))

	t.Run("Expect 1 relation after creation", func(t *testing.T) {
		rels, err := List(ctx, ListRequest{
			From: &u1,
			To:   &c1,
			Name: rel.Name,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(rels.Data))
	})
}

func TestListWithEmbeddings(t *testing.T) {
	ctx := context.Background()

	rnd := xid.New().String()
	u1 := Entity{ID: "u1" + rnd, Type: "user"}
	c1 := Entity{ID: "c1" + rnd, Type: "collection"}
	p1 := Entity{ID: "p1" + rnd, Type: "post"}

	createEntity := func(e Entity) error {
		return u.Ptr(entitiesdb.Entity{ID: e.ID, Type: e.Type}).Create(ctx)
	}

	require.NoError(t, createEntity(u1))
	require.NoError(t, createEntity(c1))
	require.NoError(t, createEntity(p1))

	u1c1, err := Create(ctx, CreateRequest{
		From: Entity{ID: u1.ID, Type: u1.Type},
		To:   Entity{ID: c1.ID, Type: c1.Type},
	})
	require.NoError(t, err)

	c1p1, err := Create(ctx, CreateRequest{
		From: Entity{ID: c1.ID, Type: c1.Type},
		To:   Entity{ID: p1.ID, Type: p1.Type},
	})
	require.NoError(t, err)

	t.Run("u1p1 depends on u1c1 and c1p1", func(t *testing.T) {
		rels, err := List(ctx, ListRequest{
			From: &u1,
			To:   &p1,
			Embed: ListEmbed{
				Dependencies: true,
			},
		})
		assert.NoError(t, err)

		assert.Equal(t, 1, len(rels.Data))
		u1p1 := rels.Data[0]
		assert.Equal(t, 2, len(rels.Data[0].DependenciesIDs))
		assert.Equal(t, u1c1.ID, u1p1.DependenciesIDs[0])
		assert.Equal(t, c1p1.ID, u1p1.DependenciesIDs[1])

		t.Run("u1c1 has u1p1 as dependant", func(t *testing.T) {
			rels, err := List(ctx, ListRequest{
				From: &u1,
				To:   &c1,
				Embed: ListEmbed{
					Dependants: true,
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(rels.Data))
			assert.Equal(t, 1, len(rels.Data[0].DependantsIDs))
			assert.Equal(t, u1p1.ID, rels.Data[0].DependantsIDs[0])
		})
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	rnd := xid.New().String()
	u1 := Entity{ID: "u1" + rnd, Type: "user"}
	c1 := Entity{ID: "c1" + rnd, Type: "collection"}
	r1 := Entity{ID: "r1" + rnd, Type: "role"}
	p1 := Entity{ID: "p1" + rnd, Type: "permission"}

	createEntity := func(e Entity) error {
		return u.Ptr(entitiesdb.Entity{ID: e.ID, Type: e.Type}).Create(ctx)
	}

	require.NoError(t, createEntity(u1))
	require.NoError(t, createEntity(c1))
	require.NoError(t, createEntity(r1))
	require.NoError(t, createEntity(p1))

	u1c1, err := Create(ctx, CreateRequest{
		From: u1,
		To:   c1,
	})
	require.NoError(t, err)

	c1r1, err := Create(ctx, CreateRequest{
		From: c1,
		To:   r1,
	})
	require.NoError(t, err)

	r1p1, err := Create(ctx, CreateRequest{
		From: r1,
		To:   p1,
	})
	require.NoError(t, err)

	relations, err := List(ctx, ListRequest{
		From: &u1,
		To:   &p1,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(relations.Data))

	var rows []string
	u.Check(pgxscan.Select(ctx, db.Conn, &rows, `
		select cache_id from dependencies where relation_id=$1
	`, u1c1.ID))

	require.NoError(t, Delete(ctx, u1c1.ID))
	require.NoError(t, Delete(ctx, c1r1.ID))
	require.NoError(t, Delete(ctx, r1p1.ID))

	relations, err = List(ctx, ListRequest{
		From: &u1,
		To:   &p1,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(relations.Data))
}

func TestCreateNamed(t *testing.T) {
	ctx := context.Background()

	rnd := xid.New().String()
	u1 := Entity{ID: "u1" + rnd, Type: "user"}
	c1 := Entity{ID: "c1" + rnd, Type: "collection"}
	r1 := Entity{ID: "r1" + rnd, Type: "role"}
	p1 := Entity{ID: "p1" + rnd, Type: "permission"}

	createEntity := func(e Entity) error {
		return u.Ptr(entitiesdb.Entity{ID: e.ID, Type: e.Type}).Create(ctx)
	}

	require.NoError(t, createEntity(u1))
	require.NoError(t, createEntity(c1))
	require.NoError(t, createEntity(r1))
	require.NoError(t, createEntity(p1))

	u1c1, err := Create(ctx, CreateRequest{
		From: u1,
		To:   c1,
	})
	require.NoError(t, err)

	c1r1, err := Create(ctx, CreateRequest{
		From: c1,
		To:   r1,
		Name: u.Ptr("foo"),
	})
	require.NoError(t, err)

	r1p1, err := Create(ctx, CreateRequest{
		From: r1,
		To:   p1,
	})
	require.NoError(t, err)

	u1p1s, err := List(ctx, ListRequest{
		From: &Entity{ID: u1.ID, Type: u1.Type},
		To:   &Entity{ID: p1.ID, Type: p1.Type},
		Name: c1r1.Name,
	})
	require.NoError(t, err)
	require.Equal(t, 1, len(u1p1s.Data))
	require.Equal(t, c1r1.Name, u1p1s.Data[0].Name)
	fmt.Println(u1c1, c1r1, r1p1)
}

func permutations[T any](ts []T) [][]T {
	nils := make([]*T, len(ts))
	return permutationsRec(ts, nils)
}

func permutationsRec[T any](ts []T, start []*T) [][]T {
	if len(ts) == 1 {
		startUnnull := make([]T, len(start))
		for i, v := range start {
			if v != nil {
				startUnnull[i] = *v
			} else {
				startUnnull[i] = ts[0]
			}
		}
		return [][]T{startUnnull}
	}

	first, rest := ts[0], ts[1:]

	perms := [][]T{}
	for i := 0; i < len(ts); i++ {
		perm := make([]*T, len(start))
		copy(perm, start)

		nilCount := 0
		for j, t := range perm {
			if t == nil {
				nilCount++
			}
			if nilCount == i+1 {
				perm[j] = &first
				break
			}
		}
		perms = append(perms, permutationsRec(rest, perm)...)
	}

	return perms
}

// TODO: benchmark
func BenchmarkF(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, err := List(ctx, ListRequest{
			From: &Entity{ID: strconv.Itoa(rand.Intn(n - 1)), Type: "user"},
			To:   &Entity{ID: strconv.Itoa(rand.Intn(n - 1)), Type: "resource"},
		})
		u.Check(err)
		// time.Sleep(time.Millisecond * 5)
	}
}
