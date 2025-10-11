-- 2_create_categories_table.up.sql
CREATE TABLE IF NOT EXISTS categories (
    id       INTEGER PRIMARY KEY,
    name     TEXT    NOT NULL,
    user_id  INTEGER NULL,
    CONSTRAINT fk_categories_user
        FOREIGN KEY (user_id) REFERENCES users(id)
);
-- Ensure (user_id, name) is unique (global categories have user_id NULL).
CREATE UNIQUE INDEX IF NOT EXISTS ux_categories_user_name ON categories(user_id, name);
-- Index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS ix_categories_user_id ON categories(user_id);
