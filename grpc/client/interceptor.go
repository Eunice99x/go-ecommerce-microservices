package client

import (
	"context"

	"github.com/eunice99x/goMicro/grpc/mapper"
	"google.golang.org/grpc"
)

// ErrorInterceptor turns gRPC status codes back into domain errors
func ErrorInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return mapper.ErrorFromStatus(invoker(ctx, method, req, reply, cc, opts...))
}
