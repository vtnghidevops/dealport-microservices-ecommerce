# base image
FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o product-service ./cmd/api

# final stage with tiny img
FROM alpine:latest

WORKDIR /app

# Install necessary tools for migrations
RUN apk --no-cache add bash

# Copy the binary from the builder stage
COPY --from=builder /app/product-service .
COPY --from=builder /app/migrations /app/migrations

# Create non-root user and set permissions
RUN adduser -D product-user && \
    chmod +x /app/product-service && \
    chown -R product-user:product-user /app

USER product-user

EXPOSE 50053

CMD ["./product-service"]

