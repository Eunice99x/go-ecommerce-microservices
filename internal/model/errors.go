package model

import "errors"

// errors shared by every layer, mapped to grpc codes and http statuses at the edges
var (
	ErrNotFound                = errors.New("not found")
	ErrAlreadyExists           = errors.New("already exists")
	ErrInvalidArgument         = errors.New("invalid argument")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
)
