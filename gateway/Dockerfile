FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY gateway/go.mod gateway/go.sum ./gateway/

WORKDIR /app/gateway

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY gateway .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o gateway ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway/gateway .
CMD ["./gateway"]
