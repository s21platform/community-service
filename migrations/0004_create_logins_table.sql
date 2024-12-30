-- +goose Up
CREATE TABLE IF NOT EXISTS logins
(
    id SERIAL PRIMARY KEY,
    login VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
)

-- +goose Down
DROP TABLE IF EXISTS logins