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
RUN CGO_ENABLED=0 GOOS=linux go build -o listener-service .

# final stage with tiny img
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/listener-service .

# Create non-root user and set permissions
RUN adduser -D listener-user && \
    chmod +x /app/listener-service && \
    chown -R listener-user:listener-user /app

USER listener-user

CMD ["./listener-service"]

