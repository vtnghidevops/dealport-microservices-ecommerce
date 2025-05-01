#!/bin/sh

# Redis initialization script
echo "Starting Redis initialization script..."

# Wait for Redis to fully start
sleep 5

# Connect to Redis and perform initialization
redis-cli << EOF
# Clear any existing data (in development only)
FLUSHALL

# Note: We no longer create sample coupons here
# Coupons will be managed through the admin interface

# Success message
SET init:status "Initialization completed successfully"
SET init:timestamp $(date +%s)
EXPIRE init:status 86400
EXPIRE init:timestamp 86400
EOF

echo "Redis initialization completed" 