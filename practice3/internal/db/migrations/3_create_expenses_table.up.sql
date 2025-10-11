-- 3_create_expenses_table.up.sql
CREATE TABLE IF NOT EXISTS expenses (
    id           INTEGER PRIMARY KEY,
    user_id      INTEGER NOT NULL,
    category_id  INTEGER NOT NULL,
    amount       NUMERIC NOT NULL,
    currency     TEXT    NOT NULL,
    spent_at     TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note         TEXT NULL,
    CONSTRAINT fk_expenses_user
        FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_expenses_category
        FOREIGN KEY (category_id) REFERENCES categories(id),
    CONSTRAINT chk_amount_positive CHECK (amount > 0)
);
-- Indexes for common queries
CREATE INDEX IF NOT EXISTS ix_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS ix_expenses_user_spent_at ON expenses(user_id, spent_at);
