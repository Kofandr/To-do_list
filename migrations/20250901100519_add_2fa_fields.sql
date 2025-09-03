-- +goose Up
ALTER TABLE users
    ADD COLUMN telegram_username VARCHAR(50) UNIQUE,
ADD COLUMN telegram_confirmed BOOLEAN DEFAULT FALSE;

CREATE TABLE twofa_codes (
                             id SERIAL PRIMARY KEY,
                             user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
                             code_hash VARCHAR(255) NOT NULL,
                             expires_at TIMESTAMP NOT NULL,
                             for_login BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_twofa_user_id ON twofa_codes(user_id);

-- +goose Down
DROP TABLE twofa_codes;
ALTER TABLE users DROP COLUMN telegram_username, DROP COLUMN telegram_confirmed;