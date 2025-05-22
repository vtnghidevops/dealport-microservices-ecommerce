FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /app

# Copy go.mod, go.sum files
COPY go.mod ./
# Copy the entire service directory
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mailService ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/mailService .
COPY --from=builder /app/templates /app/templates

# Create non-root user and set permissions
RUN adduser -D mail-user && \
    chmod +x /app/mailService && \
    chown -R mail-user:mail-user /app

USER mail-user

# Expose mail service ports
EXPOSE 50057

# Command to run the application
CMD ["./mailService"]

