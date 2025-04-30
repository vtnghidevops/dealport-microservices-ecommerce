FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cartService ./cmd/api

# Use a small alpine image for the final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cartService .

# Expose the port the app runs on
EXPOSE 50054

# Command to run the application
CMD ["/app/cartService"] 