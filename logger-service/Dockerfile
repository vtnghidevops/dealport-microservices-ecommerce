# base image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o logger-service ./cmd/api

RUN chmod +x /app/logger-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/logger-service /app

CMD ["/app/logger-service"]

