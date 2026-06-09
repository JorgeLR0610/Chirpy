-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (
    now(), now(), $1, $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteUsers :exec
TRUNCATE TABLE users;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;