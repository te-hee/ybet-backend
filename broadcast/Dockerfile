FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY broadcast/go.mod broadcast/go.sum ./broadcast/

WORKDIR /app/broadcast

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY broadcast .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o broadcast ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broadcast/broadcast .
CMD ["./broadcast"]
