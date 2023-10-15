package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/td0m/doorman/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	grpcEndpoint = flag.String("grpc-server-endpoint", "localhost:13335", "gRPC server endpoint")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEntitiesClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Create(ctx, &pb.EntitiesCreateRequest{
		Id: "foo:bar",
	})
	fmt.Println(res, err)
}
