package interceptor

import (
	"context"
	"github.com/Azzonya/go-shortener/internal/session"
	"github.com/Azzonya/go-shortener/internal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// GrpcInterceptorAuth creates a gRPC server interceptor for authentication,
// ensuring that each incoming request is authenticated using JWT tokens.
// If the user is not authenticated, a new user session is created and added
// to the outgoing context metadata with the user_id field set.
//
// Parameters:
//   - jwtSecret: Secret key used for JWT token verification.
//
// Returns:
//   - grpc.UnaryServerInterceptor: A gRPC unary server interceptor function.
func GrpcInterceptorAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		u, ok := session.GetUserFromMetadata(ctx)
		if !ok {
			var err error
			u, err = user.New()
			if err != nil {
				return nil, err
			}

			md := metadata.New(map[string]string{
				"user_id": u.ID,
			})

			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}
