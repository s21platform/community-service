-- +goose Up
CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    initiator UUID NOT NULL,
    invite_login VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS invites;
