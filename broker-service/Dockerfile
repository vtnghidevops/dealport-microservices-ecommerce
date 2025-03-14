# base image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o broker-service ./cmd/api

RUN chmod +x /app/broker-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/broker-service /app

CMD ["/app/broker-service"]

