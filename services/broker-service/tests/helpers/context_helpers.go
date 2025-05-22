package helpers

import (
	"context"
)

// Định nghĩa key type riêng biệt để tránh xung đột key trong context
type contextKey string

// Các constant cho key
const (
	userIDKey contextKey = "user_id"
)

// ContextWithUserID trả về context mới với user_id
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
