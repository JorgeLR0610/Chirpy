-- +goose Up
ALTER TABLE users ALTER COLUMN hashed_password TYPE TEXT;

-- +goose Down
ALTER TABLE users ALTER COLUMN hashed_password TYPE VARCHAR(255);
