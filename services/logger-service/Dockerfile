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
RUN CGO_ENABLED=0 GOOS=linux go build -o loggerService ./cmd/api

# final stage with tiny img
FROM alpine:latest

WORKDIR /app

# Install MongoDB tools and Node.js for running migrations
RUN apk add --no-cache mongodb-tools nodejs npm curl bash \
    python3 make g++ 

# Set PATH to include npm global binaries
ENV PATH="/usr/local/share/.config/yarn/global/node_modules/.bin:${PATH}"

# Copy the binary from the builder stage
COPY --from=builder /app/loggerService .

# Copy migration files
COPY --from=builder /app/migrations/ ./migrations/

# Make migration scripts executable
RUN chmod +x /app/migrations/*.js

EXPOSE 50056

# Default command
CMD ["./loggerService"]

