FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cart-service ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cart-service .
# Copy the migrations directory
COPY --from=builder /app/migrations /app/migrations

# Install necessary tools for Redis
RUN apk add --no-cache redis bash

# Create non-root user and set permissions
RUN adduser -D cart-user && \
    chmod +x /app/cart-service && \
    chmod +x /app/migrations/init_redis.sh && \
    chown -R cart-user:cart-user /app

USER cart-user

# Expose the port the app runs on
EXPOSE 50054

# Command to run the application
CMD ["./cart-service"] 