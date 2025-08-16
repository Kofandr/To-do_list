-- +goose Up
CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(50) NOT NULL UNIQUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       user_id INT NOT NULL,
                       completed BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_id ON tasks(user_id);


ALTER TABLE tasks
    ADD CONSTRAINT fk_user
        FOREIGN KEY (user_id) REFERENCES users(user_id)
            ON DELETE CASCADE;

-- +goose Down
ALTER TABLE tasks DROP CONSTRAINT fk_user;
DROP TABLE tasks;
DROP TABLE users;

