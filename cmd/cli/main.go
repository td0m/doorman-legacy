package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/td0m/poc-doorman/db"
	"github.com/td0m/poc-doorman/service"
)

var usage = `Usage: doorman command [options]

The command-line interface for using doorman.
Great for exploration, testing, and management.

Commands:
	entities list    lists all entities.
	entities new     creates a new entity.
	version          Prints the version.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
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

func app(ctx context.Context) error {
	entities := service.Entities{}
	relations := service.Relations{}

	cmd := os.Args[1]
	switch cmd {
	case "connectionsg":
		if len(os.Args) < 3 {
			return fmt.Errorf("subcommand required")
		}
		subcmd := os.Args[2]
		switch subcmd {
		case "list":
			// case "connect":
			// 	if len(os.Args) < 3 {
			// 		return fmt.Errorf("usage: doorman connect {{id}}")
			// 	}
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
		os.Args = os.Args[1:]
		flag.Parse()

		if *from == "" {
			from = nil
		}
		if *to == "" {
			to = nil
		}

		res, err := relations.List(ctx, service.RelationsListRequest{
			From: from,
			To:   to,
		})
		if err != nil {
			return err
		}

		for _, r := range res.Items {
			printRel(&r)
		}
	case "connect":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: doorman connections new [from] [to]")
		}

		from, to := os.Args[2], os.Args[3]

		os.Args = os.Args[3:]
		name := flag.String("name", "", "the name of the relation")
		flag.Parse()

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
		fmt.Println("👍", rel.ID)

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
		// fmt.Printf("%+v\n", entity)
	case "entities":
		if len(os.Args) < 3 {
			return fmt.Errorf("subcommand required")
		}
		subcmd := os.Args[2]
		switch subcmd {
		case "list":
		case "new":
			if len(os.Args) != 4 {
				return fmt.Errorf("usage: doorman entities new [id]")
			}
			id := os.Args[3]
			entity, err := entities.Create(ctx, service.CreateEntity{
				ID: id,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%+v\n", entity)
		}
	}

	return nil
}

func printRel(r *service.Relation) {
	fmt.Printf("%s => %s \t\t(id='%s')\n", r.From, r.To, r.ID)
}
