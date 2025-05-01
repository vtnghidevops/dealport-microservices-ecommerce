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

// CouponRepository implements domain.CouponRepository using Redis
type CouponRepository struct {
	client *redis.Client
}

// NewCouponRepository creates a new Redis-backed coupon repository
func NewCouponRepository(client *redis.Client) *CouponRepository {
	return &CouponRepository{
		client: client,
	}
}

// couponKey returns the Redis key for a coupon by ID
func couponKey(id string) string {
	return fmt.Sprintf("coupon:%s", id)
}

// couponCodeKey returns the Redis key for a coupon by code
func couponCodeKey(code string) string {
	return fmt.Sprintf("coupon:code:%s", code)
}

// couponListKey returns the Redis key for the coupon index
func couponListKey() string {
	return "coupons:list"
}

// GetCoupons retrieves a paginated list of coupons
func (r *CouponRepository) GetCoupons(ctx context.Context, page int, limit int) ([]domain.Coupon, int, error) {
	// Calculate start and end index for pagination
	start := (page - 1) * limit
	end := start + limit - 1

	// Get coupon IDs from sorted set (newest first)
	couponIDs, err := r.client.ZRevRange(ctx, couponListKey(), int64(start), int64(end)).Result()
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := r.client.ZCard(ctx, couponListKey()).Result()
	if err != nil {
		return nil, 0, err
	}

	if len(couponIDs) == 0 {
		return []domain.Coupon{}, int(total), nil
	}

	// Prepare pipeline for batch retrieval
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(couponIDs))

	for i, id := range couponIDs {
		cmds[i] = pipe.Get(ctx, couponKey(id))
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, 0, err
	}

	// Process results
	coupons := make([]domain.Coupon, 0, len(couponIDs))
	for _, cmd := range cmds {
		if cmd.Err() == redis.Nil {
			continue
		}

		data, err := cmd.Result()
		if err != nil {
			continue
		}

		var coupon domain.Coupon
		if err := json.Unmarshal([]byte(data), &coupon); err != nil {
			continue
		}

		coupons = append(coupons, coupon)
	}

	return coupons, int(total), nil
}

// GetCouponByID retrieves a coupon by its ID
func (r *CouponRepository) GetCouponByID(ctx context.Context, id string) (*domain.Coupon, error) {
	data, err := r.client.Get(ctx, couponKey(id)).Result()
	if err == redis.Nil {
		return nil, domain.ErrCouponNotFound
	} else if err != nil {
		return nil, err
	}

	var coupon domain.Coupon
	if err := json.Unmarshal([]byte(data), &coupon); err != nil {
		return nil, err
	}

	return &coupon, nil
}

// GetCouponByCode retrieves a coupon by its code
func (r *CouponRepository) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	// Get coupon ID from code index
	id, err := r.client.Get(ctx, couponCodeKey(code)).Result()
	if err == redis.Nil {
		return nil, domain.ErrCouponNotFound
	} else if err != nil {
		return nil, err
	}

	// Get coupon by ID
	return r.GetCouponByID(ctx, id)
}

// CreateCoupon creates a new coupon
func (r *CouponRepository) CreateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	// Check if code already exists
	_, err := r.client.Get(ctx, couponCodeKey(coupon.Code)).Result()
	if err != redis.Nil {
		return nil, domain.ErrCouponExists
	}

	// Generate a new ID if not provided
	if coupon.ID == "" {
		coupon.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	coupon.CreatedAt = now
	coupon.UpdatedAt = now

	// Initialize usage count
	coupon.UsageCount = 0

	// Marshal coupon to JSON
	data, err := json.Marshal(coupon)
	if err != nil {
		return nil, err
	}

	// Use a transaction to ensure data consistency
	pipe := r.client.TxPipeline()

	// Store coupon data
	pipe.Set(ctx, couponKey(coupon.ID), data, 0)

	// Create code to ID mapping
	pipe.Set(ctx, couponCodeKey(coupon.Code), coupon.ID, 0)

	// Add to sorted set for listing (sorted by creation time)
	pipe.ZAdd(ctx, couponListKey(), &redis.Z{
		Score:  float64(now.Unix()),
		Member: coupon.ID,
	})

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

// UpdateCoupon updates an existing coupon
func (r *CouponRepository) UpdateCoupon(ctx context.Context, coupon *domain.Coupon) (*domain.Coupon, error) {
	// Get existing coupon
	existing, err := r.GetCouponByID(ctx, coupon.ID)
	if err != nil {
		return nil, err
	}

	// Check if code has changed and if new code already exists
	if existing.Code != coupon.Code {
		_, err := r.client.Get(ctx, couponCodeKey(coupon.Code)).Result()
		if err != redis.Nil {
			return nil, domain.ErrCouponExists
		}
	}

	// Update timestamp
	coupon.CreatedAt = existing.CreatedAt
	coupon.UpdatedAt = time.Now()

	// Preserve usage count
	coupon.UsageCount = existing.UsageCount

	// Marshal coupon to JSON
	data, err := json.Marshal(coupon)
	if err != nil {
		return nil, err
	}

	// Use a transaction to ensure data consistency
	pipe := r.client.TxPipeline()

	// Update coupon data
	pipe.Set(ctx, couponKey(coupon.ID), data, 0)

	// Update code mappings if changed
	if existing.Code != coupon.Code {
		pipe.Del(ctx, couponCodeKey(existing.Code))
		pipe.Set(ctx, couponCodeKey(coupon.Code), coupon.ID, 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

// DeleteCoupon deletes a coupon
func (r *CouponRepository) DeleteCoupon(ctx context.Context, id string) error {
	// Get coupon to find its code
	coupon, err := r.GetCouponByID(ctx, id)
	if err != nil {
		return err
	}

	// Use a transaction to ensure data consistency
	pipe := r.client.TxPipeline()

	// Remove coupon data
	pipe.Del(ctx, couponKey(id))

	// Remove code mapping
	pipe.Del(ctx, couponCodeKey(coupon.Code))

	// Remove from list
	pipe.ZRem(ctx, couponListKey(), id)

	_, err = pipe.Exec(ctx)
	return err
}

// IncrementUsage increments the usage count for a coupon
func (r *CouponRepository) IncrementUsage(ctx context.Context, code string) error {
	// Get coupon by code
	coupon, err := r.GetCouponByCode(ctx, code)
	if err != nil {
		return err
	}

	// Increment usage count
	coupon.UsageCount++
	coupon.UpdatedAt = time.Now()

	// Marshal updated coupon
	data, err := json.Marshal(coupon)
	if err != nil {
		return err
	}

	// Update in Redis
	return r.client.Set(ctx, couponKey(coupon.ID), data, 0).Err()
}
