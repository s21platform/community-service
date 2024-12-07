-- +goose Up
CREATE TABLE IF NOT EXISTS participant
(
    id SERIAL PRIMARY KEY,
    login VARCHAR NOT NULL,
    campus_id INTEGER NOT NULL,
    class_name VARCHAR NOT NULL,
    parallel_name VARCHAR NOT NULL,
    tribe_id INT NOT NULL,
    status VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    exp_value INT,
    level INT,
    exp_to_next_level INT,
    skills JSONB,
    crp INT,
    prp INT,
    coins INT,
    badges JSONB,
    CONSTRAINT fk_participant_campus_id FOREIGN KEY (campus_id) REFERENCES campus(id),
    CONSTRAINT fk_participant_tribe_id FOREIGN KEY (tribe_id) REFERENCES tribe(id)
    );
