package relations

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	entitiesdb "github.com/td0m/poc-doorman/entities/db"
	relationsdb "github.com/td0m/poc-doorman/relations/db"
	"github.com/td0m/poc-doorman/u"
)

const n int = 1_000_000

func setupSampleData() {
	ctx := context.Background()
	//
	// fmt.Println("Creating entities...", time.Now())
	// for i := 0; i < n; i++ {
	// 	user := &entitiesdb.Entity{
	// 		Type: "user",
	// 		ID:   strconv.Itoa(i),
	// 	}
	// 	u.Check(user.Create(ctx))
	// }
	//
	// fmt.Println("Creating resources...", time.Now())
	// for i := 0; i < n; i++ {
	// 	resource := &entitiesdb.Entity{
	// 		Type: "resource",
	// 		ID:   strconv.Itoa(i),
	// 	}
	// 	u.Check(resource.Create(ctx))
	// }

	for i := 0; i < 10000; i++ {
		start := time.Now()
		params := []any{}
		query := strings.Builder{}

		query.WriteString(`
	  insert into relations(_id, from_id, from_type, to_id, to_type, attrs)
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
				map[string]any{},
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

		_, err := relationsdb.Conn.Exec(ctx, query.String(), params...)
		u.Check(err)
		fmt.Println(time.Since(start).Nanoseconds() / int64(m))
	}

	// for i := 0; i < n; i++ {
	// 	_, err := Create(ctx, CreateRequest{
	// 		From: Entity{ID: strconv.Itoa(rand.Intn(n-1)), Type: "user"},
	// 		To: Entity{ID: strconv.Itoa(rand.Intn(n-1)), Type: "resource"},
	// 	})
	// 	u.Check(err)
	// }

}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgDoorman, err := pgxpool.New(ctx, "user=doorman database=doorman")
	u.Check(err)
	defer pgDoorman.Close()

	entitiesdb.Conn = pgDoorman
	relationsdb.Conn = pgDoorman

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
			ID:    xid.New().String(),
			From:  Entity{ID: user.ID, Type: user.Type},
			To:    Entity{ID: resource.ID, Type: resource.Type},
			Attrs: map[string]any{},
		}
		relation, err := Create(ctx, req)
		require.NoError(t, err)
		require.Equal(t, req.ID, relation.ID)
		require.Equal(t, req.From.ID, relation.From.ID)
		require.Equal(t, req.From.Type, relation.From.Type)
		require.Equal(t, req.To.ID, relation.To.ID)
		require.Equal(t, req.To.Type, relation.To.Type)
		require.Equal(t, req.Attrs, relation.Attrs)
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

	t.Run("Fails on cyclic relation", func(t *testing.T) {
		// TODO
	})

	t.Run("Validates entity type", func(t *testing.T) {
		tests := []struct {
			FromType string
			ToType   string
			Success  bool
		}{
			{"collection", "collection", false},
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
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("Relating %s to %s results in sucess=%v", tt.FromType, tt.ToType, tt.Success), func(t *testing.T) {
				from := &entitiesdb.Entity{
					Type: tt.FromType,
					ID:   xid.New().String(),
				}
				require.NoError(t, from.Create(ctx))

				to := &entitiesdb.Entity{
					Type: tt.ToType,
					ID:   xid.New().String(),
				}
				require.NoError(t, to.Create(ctx))

				req := CreateRequest{
					ID:    xid.New().String(),
					From:  Entity{ID: from.ID, Type: from.Type},
					To:    Entity{ID: to.ID, Type: to.Type},
					Attrs: map[string]any{},
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

	t.Run("Sucess computing indirect relations", func(t *testing.T) {
		tests := [][]Entity{
			{{ID: "u1", Type: "user"}, {ID: "c1", Type: "collection"}, {ID: "r1", Type: "role"}, {ID: "p1", Type: "permission"}},
			// {"user", "collection", "resource"},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("Path: %+v", tt), func(t *testing.T) {
				permutations := permutations(makeRelations(tt))
				for _, pairs := range permutations {
					// Insert them from left to right
					t.Run(fmt.Sprintf("Permutation %+v", pairs), func(t *testing.T) {
						rnd := xid.New().String()
						for i := range tt {
							e := tt[i]
							tt[i] = e
							en := &entitiesdb.Entity{
								Type: e.Type,
								ID:   e.ID + ":" + rnd,
							}
							require.NoError(t, en.Create(ctx))
						}

						for _, pair := range pairs {
							req := CreateRequest{
								From: Entity{ID: pair.From.ID + ":" + rnd, Type: pair.From.Type},
								To:   Entity{ID: pair.To.ID + ":" + rnd, Type: pair.To.Type},
							}
							_, err := Create(ctx, req)
							require.NoError(t, err)
						}

						tOnly := func(t Relation) string {
							return t.From.ID + "->" + t.To.ID
						}

						for a := range tt {
							for b := range tt {
								if b <= a {
									continue
								}

								all, err := List(ctx, ListRequest{
									From: &Entity{ID: tt[a].ID + ":" + rnd, Type: tt[a].Type},
									To:   &Entity{ID: tt[b].ID + ":" + rnd, Type: tt[b].Type},
								})
								require.NoError(t, err)
								require.Equal(t, 1, len(all), "relation: %s, requests: %+v", tt[a].ID+" => "+tt[b].ID, u.Map(pairs, tOnly)) // TODO: make it 1
							}
						}
					})
				}
			})
		}
	})
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

func makeRelations(es []Entity) []Relation {
	relations := make([]Relation, len(es)-1)
	for i, left := range es {
		if i == len(es)-1 {
			break
		}
		right := es[i+1]

		relations[i] = Relation{
			From: left,
			To:   right,
		}
	}
	return relations
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
	}
}
