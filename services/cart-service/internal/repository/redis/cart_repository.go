package redis

import (
	"cart-service/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// CartRepository implements the domain.CartRepository interface using Redis
type CartRepository struct {
	client *redis.Client
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(client *redis.Client) *CartRepository {
	return &CartRepository{
		client: client,
	}
}

// GetCart retrieves a cart from Redis by user ID
func (r *CartRepository) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	// Define key for the cart in Redis
	key := fmt.Sprintf("cart:%s", userID)

	// Get the cart data from Redis
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cart not found, return a new empty cart
		cart := &domain.Cart{
			ID:     uuid.New().String(),
			UserID: userID,
			Items:  []domain.CartItem{},
			Totals: domain.CartTotals{
				Shipping: "Free",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return cart, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Unmarshal JSON data into Cart struct
	var cart domain.Cart
	err = json.Unmarshal([]byte(val), &cart)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart data: %w", err)
	}

	return &cart, nil
}

// SaveCart saves a cart to Redis
func (r *CartRepository) SaveCart(ctx context.Context, cart *domain.Cart) error {
	// Update the timestamp
	cart.UpdatedAt = time.Now()

	// Marshal cart to JSON
	cartJSON, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart data: %w", err)
	}

	// Define key for the cart in Redis
	key := fmt.Sprintf("cart:%s", cart.UserID)

	// Save to Redis with 24-hour expiration
	err = r.client.Set(ctx, key, cartJSON, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save cart: %w", err)
	}

	return nil
}

// DeleteCart removes a cart from Redis
func (r *CartRepository) DeleteCart(ctx context.Context, userID string) error {
	// Define key for the cart in Redis
	key := fmt.Sprintf("cart:%s", userID)

	// Delete from Redis
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	return nil
}
 