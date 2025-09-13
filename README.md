# Gopher Social

Gopher Social is a backend‑focused sample social media API implemented in Go (Gin + pgx). It demonstrates layered architecture (handlers → services/stores → PostgreSQL), migrations, seeding, role-based access control, and OpenAPI (Swagger) documentation. A minimal Next.js client scaffold exists but full frontend work is intentionally deferred.

## Features
- User registration & activation (token / invitation flow)
- Role-based access control (`user`, `moderator`, `admin` roles)
- Posts with tags, comments, and ownership/role checks
- Structured error responses
- Database migrations & deterministic seeding
- Swagger documentation endpoint (`/swagger/*any`)

## Tech Stack
- Go / Gin
- PostgreSQL (pgx driver & connection pool)
- Swaggo (OpenAPI generation)
- Next.js (placeholder in `web/`)

## Project Structure (excerpt)
```
.
├── internal/
│   ├── config/        # Configuration loading
│   ├── db/            # DB connection + seeding
│   ├── handler/       # HTTP handlers & middleware
│   ├── store/         # Data access layer (users, posts, comments, roles)
│   └── errors/        # Error helpers
├── migrate/migrations # SQL migration files
├── web/               # Placeholder Next.js app (frontend TODO)
├── docs/              # Generated Swagger (after running `swag init`)
└── main.go            # App entrypoint
```

## Environment Configuration
Configuration is loaded via environment variables (see `internal/config`). Common variables:

| Variable        | Description                               | Example                      |
|-----------------|-------------------------------------------|------------------------------|
| `DATABASE_URL`  | Postgres connection string                | `postgres://user:pass@localhost:5432/gopher_social?sslmode=disable` |
| `ADDR`          | Server listen address                     | `:8080`                      |
| `ENV`           | Environment (dev, prod, etc.)             | `dev`                        |
| `BASIC_AUTH_USER` / `BASIC_AUTH_PASS` | Basic auth credentials (if used) |                            |
| Auth / JWT vars | Secrets / expiration settings             | (configure as needed)        |

### Example `.env` (create manually if desired)
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gopher_social?sslmode=disable
ADDR=:8080
ENV=dev
BASIC_AUTH_USER=admin
BASIC_AUTH_PASS=changeme
```

## Getting Started
1. Install dependencies (Go modules auto-resolve on build).
2. Ensure PostgreSQL is running and accessible.
3. Export or create your `.env` values.
4. Apply migrations.
5. (Optional) Seed data.
6. Run the server.

### Migrations
```bash
make migrate-up
```
This creates the core tables (`users`, `posts`, `comments`, `roles`, etc.).

### Seeding (Optional)
```bash
go run migrate/seed/main.go
```
The seeder inserts users, posts, comments, and assigns `role_id = 1` (expects a seeded `user` role). Adjust if your role IDs differ.

### Run the Backend
```bash
make run
# or
go run .
```
Server listens on `ADDR` (default `:8080`).

## Swagger / OpenAPI
Generate docs (only needed after changing annotations):
```bash
swag init
```
Then visit: `http://localhost:8080/swagger/index.html`

## Frontend (Deferred / TODO)
The `web/` directory contains a minimal Next.js scaffold. You can ignore it for backend-focused development.
```bash
cd web
npm install
npm run dev
```
Runs at `http://localhost:3000`.

## Make Targets (if defined)
| Command            | Purpose                        |
|--------------------|--------------------------------|
| `make migrate-up`  | Apply migrations               |
| `make run`         | Run the backend server         |
| `make seed`        | (If defined) run seeding logic |

## Troubleshooting
| Symptom / Error                                        | Cause / Fix |
|--------------------------------------------------------|------------|
| `multiple default values specified` (migrations)       | Avoid `BIGSERIAL` + explicit `DEFAULT`; use `BIGINT` or separate steps. |
| `null value in column "role_id"` during seeding        | Roles not seeded yet; ensure `roles` migration ran and role id matches seed expectation. |
| `cannot scan bytea into *string`                       | Password column is `bytea`; model expects hashed bytes. Update struct types accordingly. |
| Swagger UI 500 fetching `doc.json`                     | Regenerate with `swag init`; ensure `docs` import path matches module. |
| Unauthorized or forbidden on modifying a post          | Check JWT / context user; ownership or role precedence may block. |

## Contribution Guidelines
- Prefer small, focused PRs.
- Keep SQL migrations additive (avoid destructive changes without a plan).
- Run `swag init` after adding or modifying route annotations.