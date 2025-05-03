# base image
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

# Lightweight runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mailApp .
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/.env .

# Expose mail service port
EXPOSE 9002

# Entry point for the application
CMD ["/app/mailApp"]

