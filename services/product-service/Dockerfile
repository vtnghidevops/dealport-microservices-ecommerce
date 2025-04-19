# base image
FROM golang:1.23-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o product-service ./cmd/api

RUN chmod +x /app/product-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/product-service /app
COPY data/products.sql /app/products.sql
CMD ["/app/product-service"]

