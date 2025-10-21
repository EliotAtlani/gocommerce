package middleware

import (
	"context"
	"net/http"
	"strings"

	authpb "go-project/proto/auth"
)

// contextKey is a custom type for context keys to avoid collisions
// This follows Go best practices for context.WithValue
type contextKey string

const userIDKey contextKey = "user_id"

// AuthMiddleware validates JWT tokens by calling the auth service
// This is a higher-order function (middleware pattern in Go web servers)
// It takes a handler and returns a new handler that adds authentication
func AuthMiddleware(authClient authpb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Extract and validate the token format
			// The Authorization header should be in format: "Bearer <token>"
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				http.Error(w, "Authorization header must start with 'Bearer '", http.StatusUnauthorized)
				return // CRITICAL: Must return after error response
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if token == "" {
				http.Error(w, "Token cannot be empty", http.StatusUnauthorized)
				return
			}

			// Call the auth service to validate the token
			resp, err := authClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{
				Token: token,
			})

			if err != nil {
				http.Error(w, "Failed to validate token", http.StatusInternalServerError)
				return
			}

			if !resp.Valid {
				http.Error(w, "Invalid token: "+resp.Error, http.StatusUnauthorized)
				return
			}

			// Add user_id to request context for downstream handlers to use
			// Context is Go's way of passing request-scoped values through the call chain
			ctx := context.WithValue(r.Context(), userIDKey, resp.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user_id from context in handlers
// Returns empty string if not found
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}
