-- +goose Up
Alter TABLE participant
ADD COLUMN exp_value INT,
ADD COLUMN level INT,
ADD COLUMN exp_to_next_level INT,
ADD COLUMN skills JSONB,
ADD COLUMN crp INT,
ADD COLUMN prp INT,
ADD COLUMN coins INT,
ADD COLUMN badges JSONB;

-- +goose Down
Alter TABLE participant
DROP COLUMN exp_value INT,
DROP COLUMN level INT,
DROP COLUMN exp_to_next_level INT,
DROP COLUMN crp INT,
DROP COLUMN skills JSONB,
DROP COLUMN prp INT,
DROP COLUMN coins INT,
DROP COLUMN badges JSONB;
