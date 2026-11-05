package utils

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GRPCToHTTPResponse(st *status.Status) (int,string) {
	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()

	case codes.NotFound:
		return http.StatusNotFound, st.Message()

	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()

	case codes.PermissionDenied:
		return http.StatusForbidden, st.Message()

	case codes.Unauthenticated:
		return http.StatusUnauthorized, st.Message()

	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, "Service timeout"

	case codes.Unavailable:
		return http.StatusServiceUnavailable, "Service unavailable"

	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

func SetAuth(ctx context.Context, key string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+key)
	ctxWithAuth := metadata.NewOutgoingContext(ctx, md)
	return ctxWithAuth
}
