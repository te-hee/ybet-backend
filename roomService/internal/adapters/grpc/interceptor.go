package grpcadapter

import (
	"context"
	"fmt"
	"roomService/config"
	"roomService/internal/contextkeys"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !config.Cfg.Auth.Enabled {
		return handler(ctx, req)
	}

	newCtx, err := authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return handler(newCtx, req)
}

func authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "failed to retrieve gRPC metadata")
	}

	apiKeyHeader := md.Get("x-internal-api-key")
	if len(apiKeyHeader) == 0 || apiKeyHeader[0] != config.Cfg.Auth.ServiceApiKey {
		return nil, status.Error(codes.Unauthenticated, "invalid or missing internal API key")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return contextkeys.ContextWithInternal(ctx), nil
	}

	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
	if tokenString == authHeader[0] {
		return nil, status.Error(codes.Unauthenticated, "Bearer token required in authorization header")
	}

	claims, err := verifyJWT(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid JWT: %v", err)
	}

	userUUID := claims.Subject
	if userUUID == "" {
		return nil, status.Error(codes.Unauthenticated, "user_uuid (subject) not found in JWT")
	}

	return contextkeys.ContextWithUser(ctx, userUUID, claims.Username), nil
}

func verifyJWT(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Cfg.Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
