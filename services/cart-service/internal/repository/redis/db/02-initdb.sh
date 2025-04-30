#!/bin/sh

# Redis initialization script
echo "Starting Redis initialization script..."

# Wait for Redis to fully start
sleep 5

# Connect to Redis and perform initialization
redis-cli << EOF
# Clear any existing data (in development only)
FLUSHALL

# Create sample coupons (stored as hashes)
HSET coupon:WELCOME10 discount_percent 10 min_purchase 0 expires_at 1735689600
HSET coupon:SAVE20 discount_percent 20 min_purchase 50 expires_at 1735689600
HSET coupon:FREE_SHIPPING shipping_discount 100 min_purchase 75 expires_at 1735689600

# Create coupon list for easy lookup
SADD coupons WELCOME10 SAVE20 FREE_SHIPPING

# Success message
SET init:status "Initialization completed successfully"
SET init:timestamp $(date +%s)
EXPIRE init:status 86400
EXPIRE init:timestamp 86400
EOF

echo "Redis initialization completed" 