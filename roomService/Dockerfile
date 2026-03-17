FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY roomService/go.mod roomService/go.sum ./roomService/

WORKDIR /app/roomService

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY roomService .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o roomService ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/roomService/roomService .
COPY roomService/application.yaml .
CMD ["./roomService"]
