-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, body, user_id)
VALUES (
    now(), now(), $1, $2
)
RETURNING *;