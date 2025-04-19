# base image
FROM golang:1.18-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o listener-service .

RUN chmod +x /app/listener-service

# final stage with tiny img
FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/listener-service /app

CMD ["/app/listener-service"]

