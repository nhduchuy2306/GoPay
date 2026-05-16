# gopay

Mini payment service written in Go.

## Run locally

```bash
go test ./...
go run ./cmd/server
```

## Run with Docker Compose

```bash
docker compose up --build
```

Services:
- `server` on `:8888`
- `postgres` on `:5432`
- `redis` on `:6379`
- `mailhog` SMTP on `:1025` and web UI on `:8025`

The outbox worker runs inside the server process, so there is only one application entrypoint: `cmd/server`.

## Core endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/wallets/me`
- `POST /api/v1/wallets/transfer`
- `POST /api/v1/transactions`
- `GET /api/v1/webhooks/endpoints`

## Notes

- Transaction writes are synchronous and protected by DB row locks.
- Side effects are published through the outbox pattern and processed by `cmd/worker` with bounded concurrency.
