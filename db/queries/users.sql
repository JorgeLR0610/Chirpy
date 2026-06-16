-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (
    now(), now(), $1, $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: DeleteUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateEmailAndPassword :one
UPDATE users 
SET email = $1, hashed_password = $2, updated_at = $3
WHERE id = $4
RETURNING id, email, created_at, updated_at, is_chirpy_red;

-- name: UpgradeUserToChirpyRed :execrows
UPDATE users
SET is_chirpy_red = true, updated_at = $1
WHERE id = $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;