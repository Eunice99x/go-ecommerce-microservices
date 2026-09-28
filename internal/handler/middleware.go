package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/eunice99x/goMicro/internal/pkg/auth"
)

type authKey struct{}

func GetAuthMiddlewareFunc(genToken *auth.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, genToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAdminMiddlewareFunc(genToken *auth.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := verifyClaimsFromAuthHeader(r, genToken)
			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
				return
			}

			if !claims.IsAdmin {
				http.Error(w, "user is not an admin", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func claimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(authKey{}).(*auth.Claims)
	return claims, ok && claims != nil
}

func verifyClaimsFromAuthHeader(r *http.Request, genToken *auth.JWTConfig) (*auth.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	token := fields[1]

	claims, err := genToken.ValidateToken(token, "access")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}
