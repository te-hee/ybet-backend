FROM golang:1.25.3-alpine AS builder
WORKDIR /app

COPY proto ./proto
COPY userService/go.mod userService/go.sum ./userService/

WORKDIR /app/userService

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY userService .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o userService ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/userService/userService .
COPY userService/application.yaml .
CMD ["./userService"]
