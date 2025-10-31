FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY messageService ./messageService

WORKDIR /app/messageService

RUN go build -o messageService ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/messageService/messageService .
CMD ["./messageService"]
