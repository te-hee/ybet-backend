FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY gateway ./gateway

WORKDIR /app/gateway

RUN go build -o gateway ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway/gateway .
CMD ["./gateway", "--addr=message-service:50051"]
