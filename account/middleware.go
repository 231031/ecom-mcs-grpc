package account

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md.Get("id")
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user id not found in metadata")
	}
	userID := values[0]
	newCtx := context.WithValue(ctx, "id", userID)

	values = md.Get("email")
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user email not found in metadata")
	}
	userEmail := values[0]
	newCtx = context.WithValue(newCtx, "email", userEmail)

	values = md.Get("role")
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user role not found in metadata")
	}
	userRole := values[0]
	newCtx = context.WithValue(newCtx, "role", userRole)

	return handler(newCtx, req)
}
