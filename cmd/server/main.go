package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/td0m/doorman/db"
	pb "github.com/td0m/doorman/gen/go"
	"github.com/td0m/doorman/service"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initConfigs() error {
	var validConnections [][]string
	if validConnectionsStr := os.Getenv("DOORMAN_VALID_CONNECTIONS"); len(validConnectionsStr) > 0 {
		if err := json.Unmarshal([]byte(validConnectionsStr), &validConnections); err != nil {
			return fmt.Errorf("json.Unmarshal on DOORMAN_VALID_CONNECTIONS failed: %w", err)
		}
	}

	if err := service.Setup(validConnections); err != nil {
		return fmt.Errorf("service.Setup failed: %w", err)
	}

	return nil
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := db.Init(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := initConfigs(); err != nil {
		return fmt.Errorf("initConfigs failed: %w", err)
	}

	addr := "localhost:13335"
	if envAddr := os.Getenv("DOORMAN_HOST"); len(envAddr) > 0 {
		addr = envAddr
	}

	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("net.Listen failed: %w", err)
	}

	entities := &service.Entities{}
	relations := &service.Relations{}

	s := grpc.NewServer()
	pb.RegisterEntitiesServer(s, entities)
	pb.RegisterRelationsServer(s, relations)
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
				if err := pb.RegisterEntitiesHandlerServer(ctx, mux, entities); err != nil {
					panic(fmt.Errorf("RegisterEntitiesHandlerServer failed: %w", err))
				}
				if err := pb.RegisterRelationsHandlerServer(ctx, mux, relations); err != nil {
					panic(fmt.Errorf("RegisterRelationsHandlerServer failed: %w", err))
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
