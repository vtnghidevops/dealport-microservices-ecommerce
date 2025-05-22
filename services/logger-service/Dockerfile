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
RUN CGO_ENABLED=0 GOOS=linux go build -o logger-service ./cmd/api

# final stage
FROM alpine:latest

WORKDIR /app

# Install mongodb-tools with minimal dependencies
RUN apk add --no-cache --virtual .build-deps curl bash && \
    curl -fsSL https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2004-x86_64-100.9.4.tgz | tar -xz && \
    mv mongodb-database-tools-*/bin/* /usr/local/bin/ && \
    rm -rf mongodb-database-tools-* && \
    apk del .build-deps && \
    apk add --no-cache bash && \
    rm -rf /var/cache/apk/*


# Copy migration files và binary
COPY --from=builder /app/migrations/ ./migrations/
COPY --from=builder /app/logger-service .

# Tạo non-root user và set permissions
RUN adduser -D logger-user && \
    chmod +x /app/logger-service && \
    chmod +x /app/migrations/*.js && \
    chown -R logger-user:logger-user /app

USER logger-user

EXPOSE 50056

# Default command
CMD ["./logger-service"]




