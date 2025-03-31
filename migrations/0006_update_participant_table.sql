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
DROP COLUMN exp_value,
DROP COLUMN level,
DROP COLUMN exp_to_next_level,
DROP COLUMN crp,
DROP COLUMN skills,
DROP COLUMN prp,
DROP COLUMN coins,
DROP COLUMN badges;
