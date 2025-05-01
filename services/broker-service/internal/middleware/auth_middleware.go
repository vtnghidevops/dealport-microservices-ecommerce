package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"broker-service/internal/util"
	authpb "broker-service/proto/auth"
)

type AuthMiddleware struct {
	AuthClient authpb.AuthServiceClient
}

// RequireAuth is a middleware that validates JWT tokens
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.ErrorJSON(w, errors.New("authorization header required"), http.StatusUnauthorized)
			return
		}

		// The header format should be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.ErrorJSON(w, errors.New("invalid authorization header format"), http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Call authentication service to validate token
		res, err := m.AuthClient.Validate(r.Context(), &authpb.ValidateRequest{
			Token: token,
		})

		if err != nil {
			util.ErrorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		if !res.Valid {
			util.ErrorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), "user_id", res.UserId)
		ctx = context.WithValue(ctx, "user_role", res.Claims["role"])

		// Log the user ID for debugging
		log.Printf("Auth middleware: User ID set in context: %s, Path: %s, Method: %s", res.UserId, r.URL.Path, r.Method)

		// Continue with the next handler with the enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
