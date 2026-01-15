package handlers

import (
	"context"
	"fmt"
	"log"
	"messageService/config"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if err := validate(ctx); err != nil {
		return nil, err
	}
	m, err := handler(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("RPC failed: %v", err)
	}
	return m, err
}

func validate(ctx context.Context) error {
	if !*config.NoAuth {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.DataLoss, "Failed to retrieve gRPC metadata")
		}
		authHeader := md.Get("authorization")
		log.Println(md)

		if len(authHeader) == 0 {
			return status.Errorf(codes.PermissionDenied, "auth key not provided in metadata")
		}

		key := strings.TrimPrefix(authHeader[0], "Bearer ")

		if key != config.ServiceApiKey {
			return status.Errorf(codes.Unauthenticated, "wrong auth key provided")
		}
	}
	return nil
}
