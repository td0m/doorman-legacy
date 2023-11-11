package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/server"
)

func run() error {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, "")
	if err != nil {
		return fmt.Errorf("conn failed: %w", err)
	}

	_, err = conn.Exec(ctx, `delete from tuples; delete from roles; delete from changes;`)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}

	s := server.NewDoorman(conn)

	_, _ = s.UpsertRole(ctx, &pb.UpsertRoleRequest{
		Id:    "post:viewer",
		Verbs: []string{"retrieve", "write"},
	})

	u := 1000
	p := 100

	before := time.Now()

	for i := 0; i < u; i++ {
		fmt.Println(float64(i) / float64(u))
		for j := 0; j < p; j++ {
			_, err := s.Grant(ctx, &pb.GrantRequest{
				Subject: fmt.Sprintf("user:%d", i),
				Role: "post:viewer",
				Object: fmt.Sprintf("post:%d", j),
			})
			if err != nil {
				return fmt.Errorf("fail: %w", err)
			}
		}
	}

	duration := time.Since(before).Microseconds()
	fmt.Printf("took %v. That is %d microsecs / change.\n", duration, duration / int64(u) / int64(p))

	if err := s.ProcessAllChanges(); err != nil {
		return err
	}

	duration = time.Since(before).Microseconds()
	fmt.Printf("took %v. That is %d microsecs / change.\n", duration, duration / int64(u) / int64(p))

	return nil
}


func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
