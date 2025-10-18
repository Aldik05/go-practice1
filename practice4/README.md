# Practice 4 — Postgres + sqlx (CRUD & Transaction)

## Что в проекте
- docker-compose с Postgres 16
- `users.sql` — схема и демо-данные
- Go 1.22 пример с `sqlx`: `InsertUser`, `GetAllUsers`, `GetUserByID`
- Транзакция `TransferBalance` с `SELECT ... FOR UPDATE` и защитой от гонок

## Запуск
```bash
# 1) поднять БД
docker compose up -d

# 2) зависимости
go mod tidy

# 3) запуск
go run .
