# Expense Tracker Migrations (Practice 3)

This repo contains versioned SQL migrations for a minimal Expense Tracker schema.

## Structure
```
go-practice3/
  internal/db/migrations/
    1_create_users_table.up.sql
    1_create_users_table.down.sql
    2_create_categories_table.up.sql
    2_create_categories_table.down.sql
    3_create_expenses_table.up.sql
    3_create_expenses_table.down.sql
  cmd/verify/main.go
  go.mod
```

## Choose your DB engine

- **SQLite** (easiest to start)
- **PostgreSQL** (requires a running server / Docker)

## Install golang-migrate (CLI)
See: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

On macOS (homebrew):
```
brew install golang-migrate
```

On Linux (prebuilt binaries):
Download a release for your OS/arch and `chmod +x migrate`.

## Run migrations

From the repository root:

### SQLite
```
migrate -path internal/db/migrations -database "sqlite3://./expense.db" up
migrate -path internal/db/migrations -database "sqlite3://./expense.db" down 1
migrate -path internal/db/migrations -database "sqlite3://./expense.db" version
```

### PostgreSQL
```
export PGURL="postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable"
migrate -path internal/db/migrations -database "$PGURL" up
migrate -path internal/db/migrations -database "$PGURL" down 1
migrate -path internal/db/migrations -database "$PGURL" version
```

## Notes
- `users.email` is unique, `users.name` is NOT NULL.
- `categories` may be global (`user_id` NULL) or user-specific. Composite unique `(user_id, name)` enforces per-user uniqueness (and still allows a global category with `NULL, name`).
- `expenses.amount` is constrained `> 0`. Common indexes: `user_id` and `(user_id, spent_at)`.

## Submitting
1) Push this folder to your GitHub: `github.com/<username>/go-practice3`
2) Submit the repo link.
