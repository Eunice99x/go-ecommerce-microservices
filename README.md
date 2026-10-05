# E-Commerce Microservices

An e-commerce backend built with **Go**, using a microservices architecture.

## Tech Stack

- Go
- PostgreSQL
- gRPC
- Protocol Buffers
- Docker

## Architecture

The application is composed of independent services communicating over **gRPC**, with PostgreSQL for persistent storage.

```
client ──HTTP──▶ api ──gRPC (mTLS)──▶ grpc ──▶ PostgreSQL
                                       ▲         │ notifications outbox
                     notifier ──gRPC───┘         ▼
                        └──SMTP──▶ customer email
```

- **api**: HTTP + JWT verification, no database access
- **grpc**: business logic (auth, pricing, order status) and the only service that touches the database
- **notifier**: drains the notification outbox and emails customers, with bounded concurrency

## Running it

Requires Docker and a `.env` with a `SECRET_KEY` of at least 32 characters (see `.env.example`).

```sh
make up      # dev certs, build images, start postgres + migrations + all services
make logs
make down
```

- API: http://localhost:3000
- Emails sent by the notifier: http://localhost:8025

To run the services with `go run` instead, start Postgres (`make db-up && make migrate-up`), then `make grpc-run`, `make go-run` and `make notifier-run` in separate terminals.

## API

| Method | Path | Access |
|---|---|---|
| POST | `/users` | public (sign up) |
| POST | `/login` | public |
| POST | `/refresh` | public (refresh token) |
| DELETE | `/logout` | user |
| PATCH | `/sessions/{id}/revoke` | user (own sessions) |
| GET / PATCH | `/users/user` | user (self) |
| GET | `/users` | admin |
| DELETE | `/users/{id}` | admin |
| GET | `/products`, `/products/{id}` | public |
| POST | `/products` | admin |
| PATCH / DELETE | `/products/{id}` | admin |
| POST | `/orders` | user |
| GET | `/orders/me` | user |
| GET / DELETE | `/orders/{id}` | owner or admin |
| GET | `/orders` | admin |
| PATCH | `/orders/{id}/status` | admin (emails the customer) |

## Status

Feature complete. Next: Kubernetes manifests, integration tests against a real Postgres.