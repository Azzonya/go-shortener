package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

// GrpcInterceptorLogger creates a gRPC server interceptor that logs each call,
// outputting information about the duration of execution and the response status.
func GrpcInterceptorLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		duration := time.Since(start)

		log.Printf("Request - Method: %s, Path: %s, Duration: %v", info.FullMethod, "N/A", duration)

		statusCode := status.Code(err)
		log.Printf("Response - Status: %s, Duration: %v", statusCode.String(), duration)

		return resp, err
	}
}
