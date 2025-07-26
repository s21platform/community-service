-- +goose Up
CREATE TABLE IF NOT EXISTS link_edu
(
    edu_id INTEGER NOT NULL,
    user_uuid UUID NOT NULL,
    CONSTRAINT fk_link_edu_edu_id FOREIGN KEY (edu_id) REFERENCES participant(id)
);

-- +goose Down
DROP TABLE IF EXISTS link_edu;