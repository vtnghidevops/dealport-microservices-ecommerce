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
RUN CGO_ENABLED=0 GOOS=linux go build -o productService ./cmd/api

# final stage with tiny img
FROM alpine:latest

WORKDIR /app

# Create uploads directory
RUN mkdir -p /app/uploads/products

# Copy the binary from the builder stage
COPY --from=builder /app/productService .

EXPOSE 8080
CMD ["./productService"]

