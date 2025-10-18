-- Схема и стартовые данные для демонстрации CRUD и транзакций

CREATE TABLE IF NOT EXISTS users (
    id      SERIAL PRIMARY KEY,
    name    TEXT NOT NULL,
    email   TEXT UNIQUE NOT NULL,
    balance NUMERIC(12,2) NOT NULL DEFAULT 0
);

INSERT INTO users (name, email, balance) VALUES
('Alice', 'alice@example.com', 100.00)
ON CONFLICT DO NOTHING;

INSERT INTO users (name, email, balance) VALUES
('Bob', 'bob@example.com', 50.00)
ON CONFLICT DO NOTHING;

INSERT INTO users (name, email, balance) VALUES
('Carol', 'carol@example.com', 0.00)
ON CONFLICT DO NOTHING;
