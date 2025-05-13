FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cartService ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cartService .
# Copy the migrations directory
COPY --from=builder /app/migrations /app/migrations

# Make the init script executable
RUN chmod +x /app/migrations/init_redis.sh

# Install necessary tools for Redis
RUN apk add --no-cache redis bash

# Expose the port the app runs on
EXPOSE 50054

# Command to run the application
CMD ["./cartService"] 