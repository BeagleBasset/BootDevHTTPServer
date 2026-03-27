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
