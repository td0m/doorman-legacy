package doorman

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
)

type Change struct {
	ID        string
	Type      string
	Payload   json.RawMessage
	CreatedAt time.Time
}

func TuplesToChanges(ts []Tuple, created bool) []Change {
	typ := "TUPLE_CREATED"
	if !created {
		typ = "TUPLE_REMOVED"
	}

	changes := make([]Change, len(ts))
	for i, t := range ts {
		bs, err := json.Marshal(t)
		if err != nil {
			panic(fmt.Errorf("marshal failed: %w", err))
		}
		changes[i] = Change{
			ID:      xid.New().String(),
			Type:    typ,
			Payload: bs,
		}
	}
	return changes
}

func RelationsToChanges(rs []Relation, created bool) []Change {
	typ := "RELATION_CREATED"
	if !created {
		typ = "RELATION_REMOVED"
	}

	changes := make([]Change, len(rs))
	for i, r := range rs {
		bs, err := json.Marshal(r)
		if err != nil {
			panic(fmt.Errorf("marshal failed: %w", err))
		}
		changes[i] = Change{
			ID:      xid.New().String(),
			Type:    typ,
			Payload: bs,
		}
	}
	return changes
}
