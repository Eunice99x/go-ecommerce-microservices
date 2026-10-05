package mapper

import (
	"errors"
	"log"

	"github.com/eunice99x/goMicro/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// used in both directions so server and client stay in sync
var errorCodes = []struct {
	err  error
	code codes.Code
}{
	{model.ErrNotFound, codes.NotFound},
	{model.ErrAlreadyExists, codes.AlreadyExists},
	{model.ErrInvalidArgument, codes.InvalidArgument},
	{model.ErrInvalidCredentials, codes.Unauthenticated},
	{model.ErrInvalidStatusTransition, codes.FailedPrecondition},
}

// ErrorToStatus hides unknown errors behind codes.Internal so db details don't leak
func ErrorToStatus(method string, err error) error {
	if _, ok := status.FromError(err); ok {
		return err
	}

	for _, e := range errorCodes {
		if errors.Is(err, e.err) {
			return status.Error(e.code, err.Error())
		}
	}

	log.Printf("grpc %s: %v", method, err)

	return status.Error(codes.Internal, "internal error")
}

// ErrorFromStatus lets callers use errors.Is(err, model.ErrNotFound) on gRPC errors
func ErrorFromStatus(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, e := range errorCodes {
		if st.Code() == e.code {
			return remoteError{msg: st.Message(), target: e.err}
		}
	}

	return err
}

// remoteError keeps the server's message but unwraps to the domain error
type remoteError struct {
	msg    string
	target error
}

func (e remoteError) Error() string { return e.msg }
func (e remoteError) Unwrap() error { return e.target }
