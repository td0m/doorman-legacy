package schema

import (
	"context"
	"fmt"
)

type Path []string

func (path Path) Check(ctx context.Context, r Resolver, from, to string) (bool, error) {
	if len(path) < 2 {
		return false, fmt.Errorf("bad path length")
	}

	if len(path) > 0 && path[0] == "" {
		path[0] = from
	}

	// TODO: caching to improve performance

	relations, err := r.ListForward(ctx, path[0], path[1])
	if err != nil {
		return false, fmt.Errorf("ListForward failed: %w", err)
	}

	// TODO: consider if this can be made concurrent. I think yes.
	for _, relation := range relations {
		subpath := path[1:]
		if len(subpath) == 1 {
			if relation.To == to {
				// REACHED END OF PATH! JUST RETURN
				return true, nil
			}
			continue
		}
		subpath[0] = ""

		succ, err := Path(subpath).Check(ctx, r, relation.To, to)
		if err != nil {
			return false, fmt.Errorf("path resolve failed: %w", err)
		}

		// Full path at any point means can return true
		if succ {
			return true, nil
		}
	}
	return false, nil
}
