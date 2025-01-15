-- +goose Up
CREATE TABLE IF NOT EXISTS staff
(
    id SERIAL PRIMARY KEY,
    login VARCHAR NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS staff