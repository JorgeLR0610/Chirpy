-- +goose Up
ALTER TABLE posts RENAME TO chirps;

-- +goose Down
ALTER TABLE chirps RENAME TO posts;
