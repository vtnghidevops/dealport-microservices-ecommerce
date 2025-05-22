FROM golang:1.24.3-alpine3.21 AS builder

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

# Create non-root user and set permissions
RUN adduser -D user-user && \
    chmod +x /app/user-service && \
    chown -R user-user:user-user /app

USER user-user

# Expose the port
EXPOSE 50052

# Command to run the application
CMD ["./user-service"] 