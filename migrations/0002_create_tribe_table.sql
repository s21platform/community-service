-- +goose Up
CREATE TABLE IF NOT EXISTS tribe
(
    id SERIAL PRIMARY KEY,
    campus_id INTEGER NOT NULL,
    school_tribe_id INT NOT NULL,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_tribe_campus_id FOREIGN KEY (campus_id) REFERENCES campus(id)
);

-- +goose Down
DROP TABLE IF EXISTS tribe