-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (
    now(), now(), $1, $2
)
RETURNING id, created_at, updated_at, email;

-- name: DeleteUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users 
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING id, email, created_at, updated_at;