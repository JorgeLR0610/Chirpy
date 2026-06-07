-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES (
    now(), now(), $1
)
RETURNING *;