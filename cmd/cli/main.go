package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/td0m/doorman/db"
	"github.com/td0m/doorman/service"
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

func app(ctx context.Context) error {
	entities := service.Entities{}
	relations := service.Relations{}

	cmd := os.Args[1]
	switch cmd {
	case "attrs":
		if len(os.Args) < 4 {
			usage := `usage: attrs [user:alice or #relation_id] [json]

// commands:
//   set   [key] [value]
//   unset [key] [value]
			`
			return errors.New(usage)
		}

		overwrite := flag.Bool("overwrite", false, "overwrite existing attributes instead of merging with existing values")
		os.Args = os.Args[2:]
		flag.Parse()

		_ = overwrite

		target, value := os.Args[0], flag.Arg(0)

		if target == "" {
			return errors.New("invalid user or relation")
		}

		attrs := map[string]any{}
		err := json.Unmarshal([]byte(value), &attrs)
		if err != nil {
			return fmt.Errorf("parsing value to json failed: %w", err)
		}

		if target[0] == '#' {
			return fmt.Errorf("not impl for relations yet")
		} else {
			res, err := entities.Update(ctx, service.UpdateEntity{
				ID:    target,
				Attrs: attrs,
			})
			if err != nil {
				return fmt.Errorf("update failed: %w", err)
			}

			printAttrs(res.Attrs)
		}
	case "check":
		if len(os.Args) != 4 {
			return errors.New("usage: check [from] [to]")
		}
		from, to := os.Args[2], os.Args[3]

		res, err := relations.List(ctx, service.RelationsListRequest{
			From: &from,
			To:   &to,
		})
		if err != nil {
			return err
		}

		if len(res.Items) > 0 {
			fmt.Println("✅")
		} else {
			fmt.Println("🚫")
			os.Exit(1)
		}

	case "connections":
		from := flag.String("from", "", "from")
		to := flag.String("to", "", "to")
		name := flag.String("name", "", "name")

		os.Args = os.Args[1:]
		flag.Parse()

		if *from == "" {
			from = nil
		}
		if *to == "" {
			to = nil
		}
		if *name == "" {
			name = nil
		}

		if from == nil && to == nil {
			return errors.New("neither -from nor -to specified")
		}


		res, err := relations.List(ctx, service.RelationsListRequest{
			From: from,
			To:   to,
			Name: name,
		})
		if err != nil {
			return err
		}

		printRels(res.Items)
	case "connect":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: doorman connect [from] [to]")
		}
		os.Args = os.Args[1:]
		name := flag.String("name", "", "the name of the relation")
		flag.Parse()

		from, to := flag.Arg(0), flag.Arg(1)

		if *name == "" {
			name = nil
		}

		rel, err := relations.Create(ctx, service.RelationsCreate{
			From: from,
			To:   to,
			Name: name,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if name == nil {
			empty := ""
			name = &empty
		}
		fmt.Println("👍", rel.ID, *name)

	case "new":
		if len(os.Args) != 3 {
			return fmt.Errorf("usage: doorman new [id]")
		}
		id := os.Args[2]
		_, err := entities.Create(ctx, service.CreateEntity{
			ID: id,
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("👍")
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

func printRels(rs []service.Relation) {

	rows := [][]string{}
	for _, r := range rs {
		id := " - "
		if !strings.HasPrefix(r.ID, "cache:") {
			id = r.ID
		}

		name := ""
		if r.Name != nil {
			name = *r.Name
		}
		rows = append(rows, []string{id, emojify(r.From), name, emojify(r.To)})
	}
	table := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("ID", "From", "Name", "To").
    StyleFunc(func(row, _ int) lipgloss.Style {
        switch row {
        case 0:
            return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Padding(0, 2)
        default:
            return lipgloss.NewStyle().Padding(0, 1)
        }
    }).
		Rows(rows...)

	fmt.Println(table.Render())
}

func main() {
	usage = strings.Replace(usage, "{{version}}", "v0", 1)
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	ctx := context.Background()

	err := db.Init(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := app(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
