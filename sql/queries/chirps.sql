-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(), -- automatically generate UUID
    NOW(),             -- created_at timestamp
    NOW(),             -- updated_at timestamp
    $1,                -- body parameter
    $2                 -- user_id parameter
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetChirpFromUser :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at ASC;
