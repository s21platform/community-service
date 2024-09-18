-- +goose Up
CREATE TABLE IF NOT EXISTS campus
(
    id        SERIAL PRIMARY KEY,
    campus_uuid UUID NOT NULL,
    short_name      VARCHAR(255),
    full_name      VARCHAR(255),
    created_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS campus;