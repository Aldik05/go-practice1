-- 3_create_expenses_table.down.sql
DROP INDEX IF EXISTS ix_expenses_user_spent_at;
DROP INDEX IF EXISTS ix_expenses_user_id;
DROP TABLE IF EXISTS expenses;
