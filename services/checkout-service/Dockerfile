FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required packages for postgres
RUN apk add --no-cache gcc musl-dev

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o checkoutService ./cmd/api

# Use a small alpine image for the final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/checkoutService .

# Copy migrations
COPY --from=builder /app/scripts/migrations ./scripts/migrations

# Expose the port the app runs on
EXPOSE 50055

# Command to run the application
CMD ["/app/checkoutService"] 