FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/user-service .
# Copy the migrations directory
COPY --from=builder /app/migrations /app/migrations

# Install necessary tools for migrations
RUN apk add --no-cache postgresql-client bash

# Expose the port
EXPOSE 50052

# Command to run the application
CMD ["./user-service"] 