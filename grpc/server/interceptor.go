package server

import (
	"context"

	"github.com/eunice99x/goMicro/grpc/mapper"
	"google.golang.org/grpc"
)

// ErrorInterceptor converts every handler error into a gRPC status code
func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, mapper.ErrorToStatus(info.FullMethod, err)
	}

	return resp, nil
}
