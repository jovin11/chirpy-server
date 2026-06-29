# Chirpy Server

Chirpy is a RESTful microblogging API built with Go and PostgreSQL. It supports user accounts, secure authentication, short-form posts ("chirps"), author-based filtering, chronological sorting, and an authenticated webhook for account upgrades.

This project is an exploration of building high-performance, RESTful HTTP services with Go's standard library, focusing on resource-based routing, middleware, scalable API design, type-safe database access, and explicit authentication state while leveraging Go's concurrency model to build standardized, maintainable web servers.

## Features

- Create users and update account credentials
- Hash passwords with Argon2id
- Authenticate with short-lived JWT access tokens
- Issue, refresh, and revoke long-lived refresh tokens
- Create, retrieve, filter, sort, and delete chirps
- Enforce a 140-character chirp limit and replace selected disallowed words
- Restrict chirp deletion to its author
- Process API-key-authenticated Chirpy Red upgrade webhooks
- Serve static files and track file-server visits with atomic middleware
- Expose health, metrics, and development-only reset endpoints
- Persist data in PostgreSQL using sqlc-generated Go code

## Tech stack

- Go 1.25
- PostgreSQL
- `net/http` with method-aware routing
- sqlc for type-safe database access
- Goose-style SQL migrations
- JWT (HS256) and Argon2id

## Getting started

### Prerequisites

- Go 1.25 or newer
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for migrations
- [sqlc](https://sqlc.dev/) only if you want to regenerate database code

### Configuration

Create a PostgreSQL database, then add a `.env` file at the repository root:

```env
DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=replace-with-a-long-random-secret
POLKA_KEY=replace-with-your-webhook-api-key
```

`PLATFORM=dev` enables the database reset endpoint. Use another value outside local development.

### Migrate and run

Apply the migrations using the same connection string configured in `DB_URL`:

```bash
goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
go run .
```

The server listens on `http://localhost:8080`. Check that it is ready with:

```bash
curl http://localhost:8080/api/healthz
```

## API

Protected endpoints expect an access token in this form:

```http
Authorization: Bearer <access-token>
```

| Method | Endpoint | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | None | Readiness check |
| `POST` | `/api/users` | None | Create a user |
| `PUT` | `/api/users` | Bearer JWT | Update the current user's email and password |
| `POST` | `/api/login` | None | Return an access token and refresh token |
| `POST` | `/api/refresh` | Bearer refresh token | Issue a new access token |
| `POST` | `/api/revoke` | Bearer refresh token | Revoke a refresh token |
| `POST` | `/api/chirps` | Bearer JWT | Create a chirp |
| `GET` | `/api/chirps` | None | List chirps, optionally filtered and sorted |
| `GET` | `/api/chirps/{chirpID}` | None | Get one chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Bearer JWT | Delete a chirp owned by the current user |
| `POST` | `/api/polka/webhooks` | API key for upgrade events | Process a Chirpy Red upgrade event |
| `GET` | `/admin/metrics` | None | Show static file-server visit count |
| `POST` | `/admin/reset` | Development only | Delete application data and reset metrics |

### Chirp query parameters

`GET /api/chirps` supports:

- `author_id=<uuid>` to return chirps from one user
- `sort=asc` for oldest-first ordering
- `sort=desc` for newest-first ordering (the default)

The parameters can be combined:

```bash
curl "http://localhost:8080/api/chirps?author_id=<user-id>&sort=asc"
```

### Example flow

Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"bird@example.com","password":"correct-horse-battery-staple"}'
```

Log in:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bird@example.com","password":"correct-horse-battery-staple"}'
```

Create a chirp using the returned access token:

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer <access-token>" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello from Chirpy!"}'
```

### Webhook authentication

Chirpy Red upgrade events use an API key:

```http
Authorization: ApiKey <POLKA_KEY>
```

Example payload:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<user-id>"
  }
}
```

## Development

Run the test suite:

```bash
go test ./...
```

After changing SQL schemas or queries, regenerate the database package:

```bash
sqlc generate
```

Database migrations live in `sql/schema`, handwritten queries in `sql/queries`, and generated database code in `internal/database`.

## Project structure

```text
.
├── internal/
│   ├── auth/          # Password, JWT, refresh-token, and API-key helpers
│   └── database/      # sqlc-generated models and queries
├── sql/
│   ├── queries/       # Application queries
│   └── schema/        # Ordered database migrations
├── handler_*.go       # HTTP handlers and middleware
├── main.go            # Configuration, routing, and server startup
└── sqlc.yaml          # sqlc configuration
```
