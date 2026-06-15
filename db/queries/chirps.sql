-- name: CreateChirp :one
INSERT INTO chirps (created_at, updated_at, body, user_id)
VALUES (
    now(), now(), $1, $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :execrows
DELETE FROM chirps WHERE id = $1;