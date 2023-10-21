package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var usage = `doorman {{version}}
  The official command-line interface for doorman.

usage:
  doorman command [options]

commands:
  connect       connect an entity to another one.
  disconnect    remove a connection.
  check
  info
	list-backward
`

var (
	c pb.DoormanClient
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

	c = pb.NewDoormanClient(conn)

	cmd := os.Args[1]
	switch cmd {

	case "check":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: doorman check [from] [name] [to]")
		}

		os.Args = os.Args[1:]
		flag.Parse()

		from, name, to := flag.Arg(0), flag.Arg(1), flag.Arg(2)

		r, err := c.Check(ctx, &pb.CheckRequest{
			From: from,
			To:   to,
			Name: name,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", r)

	case "connect":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: doorman connect [from] [name] [to]")
		}

		os.Args = os.Args[1:]
		flag.Parse()

		from, name, to := flag.Arg(0), flag.Arg(1), flag.Arg(2)

		r, err := c.Connect(ctx, &pb.ConnectRequest{
			From: from,
			To:   to,
			Name: name,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", r)

	case "disconnect":
		if len(os.Args) < 4 {
			return fmt.Errorf("usage: doorman disconnect [from] [name] [to]")
		}

		os.Args = os.Args[1:]
		flag.Parse()

		from, name, to := flag.Arg(0), flag.Arg(1), flag.Arg(2)

		r, err := c.Disconnect(ctx, &pb.DisconnectRequest{
			From: from,
			To:   to,
			Name: name,
		})
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", r)
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

func main() {
	usage = strings.Replace(usage, "{{version}}", "v0", 1)
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	if err := app(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
