-- +goose Up
CREATE TABLE IF NOT EXISTS login
(
    id SERIAL PRIMARY KEY,
    nickname VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
    );

-- +goose Down
DROP TABLE IF EXISTS login