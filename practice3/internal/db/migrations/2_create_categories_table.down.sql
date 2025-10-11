-- 2_create_categories_table.down.sql
DROP INDEX IF EXISTS ix_categories_user_id;
DROP INDEX IF EXISTS ux_categories_user_name;
DROP TABLE IF EXISTS categories;
