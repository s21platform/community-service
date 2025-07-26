-- +goose Up
CREATE TABLE IF NOT EXISTS pending_link
(
    participant_login TEXT NOT NULL,
    code INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS pending_link