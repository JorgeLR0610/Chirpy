-- +goose Up
CREATE TABLE posts(
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp not null,
    body varchar not null,
    user_id UUID,

    constraint user_fk
        FOREIGN KEY (user_id)
        references users(id)
        on delete cascade
);

-- +goose Down
DROP TABLE posts;