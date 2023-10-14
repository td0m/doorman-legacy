package service

import (
	"context"
	"fmt"

	"github.com/td0m/poc-doorman/db"
)

func CreateRelation(ctx context.Context) error {
	r := &db.Relation{}
	if err := r.Create(ctx); err != nil {
		return fmt.Errorf("create failed: %w", err)
	}
	return nil
}
