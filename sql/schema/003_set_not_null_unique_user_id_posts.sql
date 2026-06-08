-- +goose Up
ALTER TABLE posts ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
ALTER TABLE posts ALTER COLUMN user_id DROP NOT NULL;