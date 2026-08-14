# fark-tee-backend

Backend service for fark-tee, built with Go.

## Stack

- [huma](https://github.com/danielgtaylor/huma) + [echo](https://github.com/labstack/echo) - HTTP API
- [mongo-driver](https://github.com/mongodb/mongo-go-driver) - MongoDB
- [wire](https://github.com/google/wire) - dependency injection
- [mockery](https://github.com/vektra/mockery) - mock generation

## Prerequisites

- Go 1.26+
- Docker (for local MongoDB and S3-compatible storage)

This module depends on a private GitHub module, so configure `GOPRIVATE` before fetching dependencies:

```sh
go env -w GOPRIVATE="$(go env GOPRIVATE),github.com/fark-tee"
```

## Getting started

1. Copy the environment file and adjust as needed:

   ```sh
   cp .env.example .env
   ```

2. Start MongoDB and the S3-compatible storage (MinIO):

   ```sh
   docker compose up -d
   ```

3. Run the server:

   ```sh
   make start
   ```

## Development

Common tasks are wired up via the `Makefile`:

| Command | Description |
| --- | --- |
| `make start` | Run the server |
| `make build` | Build binaries into `bin/` |
| `make generate` | Regenerate mocks, wire DI, then tidy and format |
| `make mock` | Regenerate mocks with mockery |
| `make wire` | Regenerate wire DI code |
| `make fmt` | Format code |
| `make tidy` | Tidy Go modules |
| `make tools-install` | Install mockery |

## Project layout

```
cmd/            entrypoints (server)
internal/       application code (config, infrastructure, router, ...)
pkg/            shared packages exposed outside internal use
```
