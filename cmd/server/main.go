package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/td0m/poc-doorman/gen/go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

type Entities struct {
	*pb.UnimplementedEntitiesServiceServer
}


func (e *Entities) Create(ctx context.Context, request *pb.EntitiesCreateRequest) (*pb.Entity, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	entities := &Entities{}

	sock, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		return fmt.Errorf("net.Listen failed: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterEntitiesServiceServer(s, entities)
	reflection.Register(s)

	// Capture os signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func(sock net.Listener) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				s.ServeHTTP(w, r)
			} else {
				mux := runtime.NewServeMux()
				err := pb.RegisterEntitiesServiceHandlerServer(ctx, mux, entities)
				fmt.Println(err)
				mux.ServeHTTP(w, r)
			}
		})

		if err := http.Serve(sock, h2c.NewHandler(handler, &http2.Server{})); err != nil {
			panic(fmt.Errorf("http.Serve failed: %w", err))
		}

	}(sock)

	<-sigchan
  return nil
}

func main() {

  if err := run(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
