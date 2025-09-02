-- ./migrations/20250901120000_add_telegram_fields.sql
-- +goose Up
ALTER TABLE users
    ADD COLUMN telegram_chat_id BIGINT,
ADD COLUMN link_code TEXT,
ADD COLUMN telegram_code TEXT,
ADD COLUMN code_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN telegram_chat_id,
DROP COLUMN link_code,
DROP COLUMN telegram_code,
DROP COLUMN code_expires_at;
