package doorman

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrTupleExists   = status.Error(codes.AlreadyExists, "tuple already exists")
	ErrTupleNotFound = status.Error(codes.NotFound, "tuple not found")
)
