package contextkeys

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	userUUIDKey   contextKey = "user_uuid"
	usernameKey   contextKey = "username"
	isInternalKey contextKey = "is_internal"
)

func ContextWithUser(ctx context.Context, userUUID string, username string) context.Context {
	ctx = context.WithValue(ctx, userUUIDKey, userUUID)
	ctx = context.WithValue(ctx, usernameKey, username)
	return ctx
}

func ContextWithInternal(ctx context.Context) context.Context {
	return context.WithValue(ctx, isInternalKey, true)
}

func UserUUIDFromContext(ctx context.Context) (string, error) {
	v, ok := ctx.Value(userUUIDKey).(string)
	if !ok || v == "" {
		return "", status.Error(codes.Unauthenticated, "user_uuid not found in context")
	}
	return v, nil
}

func UsernameFromContext(ctx context.Context) (string, error) {
	v, ok := ctx.Value(usernameKey).(string)
	if !ok || v == "" {
		return "", status.Error(codes.Unauthenticated, "username not found in context")
	}
	return v, nil
}

func IsInternalRequest(ctx context.Context) bool {
	v, _ := ctx.Value(isInternalKey).(bool)
	return v
}
