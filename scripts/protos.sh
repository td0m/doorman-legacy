#!/usr/bin/env bash

protoc -I include/googleapis -I proto --go_out=./gen/go --go-grpc_out=./gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		doorman.proto

mkdir -p gen/go

protoc -I include/googleapis -I proto --grpc-gateway_out ./gen/go \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		doorman.proto

