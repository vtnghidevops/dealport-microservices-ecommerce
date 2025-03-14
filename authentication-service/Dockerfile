# base image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o authentication-service ./cmd/api

RUN chmod +x /app/authentication-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/authentication-service /app

CMD ["/app/authentication-service"]

