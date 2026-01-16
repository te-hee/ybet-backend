package auth

import (
	"broadcast/config"
	"context"
	"fmt"
	"gateway/internal/model"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserIDKey = "userID"
const UsernameKey = "username"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if *config.NoAuth {
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

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func verifyJwt(tokenString string) (*model.UserClaims, error) {
	var claims model.UserClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*model.UserClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unsupported claim type")
	}
}
