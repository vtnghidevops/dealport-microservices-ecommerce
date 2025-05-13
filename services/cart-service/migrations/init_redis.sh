#!/bin/sh

echo "Starting Redis initialization script..."

# Chờ Redis khởi động (đổi từ 127.0.0.1 sang redis-cart)
until redis-cli -h redis-cart ping | grep -q PONG; do
  echo "Waiting for Redis at redis-cart..."
  sleep 1
done

redis-cli -h redis-cart << EOF
FLUSHALL
SET init:status "Initialization completed successfully"
SET init:timestamp $(date +%s)
EXPIRE init:status 86400
EXPIRE init:timestamp 86400
EOF

echo "Redis initialization completed"
