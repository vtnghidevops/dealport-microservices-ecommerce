FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o authentication-service ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/authentication-service .
# Copy the migrations directory
COPY --from=builder /app/migrations /app/migrations

# Install necessary tools for migrations
RUN apk add --no-cache postgresql-client bash

# Create non-root user and set permissions
RUN adduser -D authentication-user && \
    chmod +x /app/authentication-service && \
    chown -R authentication-user:authentication-user /app

USER authentication-user

EXPOSE 50051
# Command to run the application
CMD ["./authentication-service"] 