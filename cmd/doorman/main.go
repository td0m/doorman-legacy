package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var usage = `doorman {{version}}
  The official command-line interface for doorman.

usage:
  doorman command [options]

commands:
  attrs        updates attributes of an entity or collection.
  connect      connect an entity to another one.
  delete       delete an existing entity.
  disconnect   remove a connection.
  new          create an entity.
`

var (
	srv pb.DoormanClient
)

func app(ctx context.Context) error {
	addr := "localhost:13335"
	if envAddr := os.Getenv("DOORMAN_HOST"); len(envAddr) > 0 {
		addr = envAddr
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc.Dial failed: %w", err)
	}
	defer conn.Close()

	srv = pb.NewDoormanClient(conn)

	cmd := os.Args[1]
	switch cmd {
	// case "attrs":
	// 	if len(os.Args) < 4 {
	// 		usage := `usage: attrs [user:alice or #relation_id] [json]`
	// 		return errors.New(usage)
	// 	}
	//
	// 	overwrite := flag.Bool("overwrite", false, "overwrite existing attributes instead of merging with existing values")
	// 	os.Args = os.Args[2:]
	// 	flag.Parse()
	//
	// 	_ = overwrite
	//
	// 	target, value := os.Args[0], flag.Arg(0)
	//
	// 	if target == "" {
	// 		return errors.New("invalid user or relation")
	// 	}
	//
	// 	attrs := map[string]any{}
	// 	err := json.Unmarshal([]byte(value), &attrs)
	// 	if err != nil {
	// 		return fmt.Errorf("parsing value to json failed: %w", err)
	// 	}
	//
	// 	pbAttrs, err := structpb.NewStruct(attrs)
	// 	if err != nil {
	// 		return fmt.Errorf("structpb.NewStruct failed: %w", err)
	// 	}
	//
	// 	if target[0] == '#' {
	// 		return fmt.Errorf("not impl for relations yet")
	// 	} else {
	// 		res, err := entities.Update(ctx, &pb.EntitiesUpdateRequest{
	// 			Id:    target,
	// 			Attrs: pbAttrs,
	// 		})
	// 		if err != nil {
	// 			return fmt.Errorf("update failed: %w", err)
	// 		}
	//
	// 		printAttrs(res.Attrs.AsMap())
	// 	}
	case "check":
		if len(os.Args) != 5 {
			return errors.New("usage: check [subject] [verb] [object]")
		}
		subject, verb, object := os.Args[2], os.Args[3], os.Args[4]

		res, err := srv.Check(ctx, &pb.CheckRequest{
			Subject: subject,
			Verb:    verb,
			Object:  object,
		})
		if err != nil {
			return err
		}

		if res.Success {
			fmt.Println("✅")
		} else {
			fmt.Println("X")
		}

	case "revoke":
		if len(os.Args) != 5 {
			return errors.New("usage: revoke [subject] [role] [object]")
		}
		subject, role, object := os.Args[2], os.Args[3], os.Args[4]

		res, err := srv.Revoke(ctx, &pb.RevokeRequest{
			Subject: subject,
			Role:    role,
			Object:  object,
		})
		if err != nil {
			return err
		}
		fmt.Println(res)

	case "grant":
		if len(os.Args) != 5 {
			return errors.New("usage: grant [subject] [role] [object]")
		}
		subject, role, object := os.Args[2], os.Args[3], os.Args[4]

		res, err := srv.Grant(ctx, &pb.GrantRequest{
			Subject: subject,
			Role:    role,
			Object:  object,
		})
		if err != nil {
			return err
		}
		fmt.Println(res)

	case "roles":
		os.Args = os.Args[1:]
		switch os.Args[1] {
		case "upsert":
			if len(os.Args) < 2 {
				return errors.New("usage: roles create [id] [verb1] ... [verbN]")
			}
			id, verbs := os.Args[2], os.Args[3:]

			role, err := srv.UpsertRole(ctx, &pb.UpsertRoleRequest{
				Id:    id,
				Verbs: verbs,
			})
			if err != nil {
				return fmt.Errorf("upsert failed: %w", err)
			}
			fmt.Println(role)
		default:
			return fmt.Errorf("invalid command: %s", os.Args[1])
		}
	}

	return nil
}

func emojify(id string) string {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 1 {
		return id
	}
	rest := parts[1]
	if parts[0] == "user" {
		return "U " + rest
	}
	return id
}

func printAttrs(attrs map[string]any) {
	for k, v := range attrs {
		fmt.Printf("%s\t\t%+v\n", k, v)
	}
}

// func printRels(rs []*pb.Relation) {
//
//		rows := [][]string{}
//		for _, r := range rs {
//			id := " - "
//			if !strings.HasPrefix(r.Id, "cache:") {
//				id = r.Id
//			}
//
//			name := ""
//			if r.Name != nil {
//				name = *r.Name
//			}
//			rows = append(rows, []string{id, emojify(r.FromId), name, emojify(r.ToId)})
//		}
//		table := table.New().
//			Border(lipgloss.NormalBorder()).
//			Headers("ID", "From", "Name", "To").
//			StyleFunc(func(row, _ int) lipgloss.Style {
//				switch row {
//				case 0:
//					return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Padding(0, 2)
//				default:
//					return lipgloss.NewStyle().Padding(0, 1)
//				}
//			}).
//			Rows(rows...)
//
//		fmt.Println(table.Render())
//	}
func main() {
	usage = strings.Replace(usage, "{{version}}", "v0", 1)
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	if err := app(ctx); err != nil {
		if st, ok := status.FromError(err); ok {
			fmt.Printf("(%s) %s\n", st.Code(), st.Message())
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}
