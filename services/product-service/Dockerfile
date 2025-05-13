# base image
FROM golang:1.24.1-alpine AS builder

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

# Create uploads directory
RUN mkdir -p /app/uploads/products

# Install necessary tools for migrations
RUN apk --no-cache add bash

# Copy the binary from the builder stage
COPY --from=builder /app/product-service .
COPY --from=builder /app/uploads/products /app/uploads/products
COPY --from=builder /app/migrations /app/migrations

EXPOSE 50053

CMD ["./product-service"]

