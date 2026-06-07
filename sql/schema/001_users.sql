-- +goose Up
CREATE TABLE users(
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp not null,
    email varchar unique not null
);

-- +goose Down
DROP TABLE users;