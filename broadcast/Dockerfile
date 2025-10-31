FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY broadcast ./broadcast

WORKDIR /app/broadcast

RUN go build -o broadcast ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broadcast/broadcast .
CMD ["./broadcast", "--addr=message-service:50051"]
