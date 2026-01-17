FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY messageService/go.mod messageService/go.sum ./messageService/

WORKDIR /app/messageService

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY messageService .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o messageService ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/messageService/messageService .
CMD ["./messageService"]
