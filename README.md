# Gopher Social

Gopher Social is a small social media sample application built in Go. It includes a REST API backend (Gin + pgx) and a Next.js frontend.

## Quickstart

Requirements
- Go 1.20+
- PostgreSQL
- Node.js 18+
- Make

### 1) Configure
Copy or edit `internal/config` environment variables or use the provided `.env` if available. Configure your `DATABASE_URL`, `ADDR`, and other settings.

### 2) Run migrations
Apply DB migrations (uses `migrate` in Makefile):

```bash
make migrate-up
```

This will create the `users`, `posts`, `roles`, and other tables. If you run into `multiple default values` errors, ensure migrations were updated to not use `BIGSERIAL` with explicit `DEFAULT` values.

### 3) Seed the database (optional)
To populate the DB with sample data:

```bash
go run migrate/seed/main.go
```

If you encounter `null value in column "role_id"` errors while seeding, ensure migrations have been applied and that the `roles` table contains a `user` role. The seeder sets `role_id=1` for seeded users by default.

### 4) Run the backend

```bash
make run
# or
go run .
```

The server will listen on the address configured in `internal/config` (default `:8080`).

### 5) Run the frontend

```bash
cd web
npm install
npm run dev
```

The Next.js app runs on `http://localhost:3000` by default.

## API docs (Swagger)
If you run `swag init` in the project root the `docs` package will be generated. The Swagger UI is mounted at `/swagger/*any`.

## Troubleshooting
- "multiple default values specified for column": edit the migration to avoid adding `BIGSERIAL` with a `DEFAULT` value. Use `BIGINT DEFAULT ...` or add column without default.
- Seeder fails with `role_id` NOT NULL: ensure roles are created and the seeder's `RoleID` matches the intended role id.
- Password scanning errors: the `users.password` column is `bytea` and the code expects a hashed `[]byte`.