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
			log.Printf("AUTH ERROR: Token validation failed: %v", err)
			util.ErrorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		if !res.Valid {
			log.Printf("AUTH ERROR: Invalid token for request: %s", r.URL.Path)
			util.ErrorJSON(w, errors.New("invalid token"), http.StatusUnauthorized)
			return
		}

		// Extract role from claims
		userRole := res.Claims["role"]

		// Add user information to request context
		ctx := context.WithValue(r.Context(), "user_id", res.UserId)
		ctx = context.WithValue(ctx, "role", userRole) // Note: using "role" key consistently

		// Log the user ID and role for debugging
		log.Printf("AUTH: User authenticated - ID: %s, Role: %s, Path: %s, Method: %s",
			res.UserId, userRole, r.URL.Path, r.Method)

		// Continue with the next handler with the enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin is a middleware that ensures the user has admin role
func (m *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First check if user is authenticated
		role, ok := r.Context().Value("role").(string)
		if !ok {
			log.Printf("ADMIN CHECK: No role found in context")
			util.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		// Check if role is admin
		if role != "admin" {
			log.Printf("ADMIN CHECK: User has role '%s', admin required", role)
			util.ErrorJSON(w, errors.New("admin privileges required"), http.StatusForbidden)
			return
		}

		log.Printf("ADMIN CHECK: Admin access granted for path: %s", r.URL.Path)

		// User is admin, continue
		next.ServeHTTP(w, r)
	})
}
