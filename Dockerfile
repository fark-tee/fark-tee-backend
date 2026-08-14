# syntax=docker/dockerfile:1

# --- Stage 1: Builder ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ENV GOPRIVATE=github.com/fark-tee/*

COPY go.mod go.sum ./

RUN --mount=type=secret,id=github_token,env=GITHUB_TOKEN \
    git config --global url."https://x-access-token:${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/" && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/server ./cmd/server

# --- Stage 2: Runtime ---
FROM alpine:3.22 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates

ENV PORT=8080

RUN addgroup -g 10001 app && \
    adduser -D -u 10001 -G app app

COPY --from=builder --chown=app:app /out/server ./server

USER 10001

EXPOSE 8080

CMD ["./server"]
