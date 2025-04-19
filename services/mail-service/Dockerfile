# base image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o mail-service ./cmd/api

RUN chmod +x /app/mail-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/mail-service /app
COPY --from=builder /app/templates /app/templates


CMD ["/app/mail-service"]

