package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	addr := "localhost:13335"
	if envAddr := os.Getenv("DOORMAN_HOST"); len(envAddr) > 0 {
		addr = envAddr
	}

	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("net.Listen failed: %w", err)
	}

	conn, err := pgxpool.New(ctx, "")
	if err != nil {
		return fmt.Errorf("pgxpool.New failed: %w", err)
	}

	srv := server.NewDoorman(conn)

	s := grpc.NewServer()
	pb.RegisterDoormanServer(s, srv)
	reflection.Register(s)

	// Capture os signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting server on: %s\n", addr)

	go func(sock net.Listener) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
				s.ServeHTTP(w, r)
			} else {
				mux := runtime.NewServeMux()
				if err := pb.RegisterDoormanHandlerServer(ctx, mux, srv); err != nil {
					panic(fmt.Errorf("RegisterDoormanHandlerServer failed: %w", err))
				}
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
