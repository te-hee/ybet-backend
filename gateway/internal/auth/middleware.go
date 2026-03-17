package auth

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/model"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserIDKey = "userID"
const UsernameKey = "username"
const RawTokenKey = "rawToken"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !config.Cfg.Auth.Enabled {
			ctx := context.WithValue(r.Context(), UserIDKey, "8f0d8552-d07d-432d-9018-8374313f9151")
			ctx = context.WithValue(ctx, UsernameKey, "cutiepie")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header not provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		claims, err := verifyJwt(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserId)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, RawTokenKey, tokenString)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func verifyJwt(tokenString string) (*model.UserClaims, error) {
	var claims model.UserClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(config.Cfg.Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*model.UserClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unsupported claim type")
	}
}

func UserFromContext(ctx context.Context) (string, string) {
	userID, _ := ctx.Value(UserIDKey).(string)
	username, _ := ctx.Value(UsernameKey).(string)
	return userID, username
}

func RawTokenFromContext(ctx context.Context) string {
	token, _ := ctx.Value(RawTokenKey).(string)
	return token
}
