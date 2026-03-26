-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), -- automatically generate UUID
    NOW(),             -- created_at timestamp
    NOW(),             -- updated_at timestamp
    $1                 -- email parameter
)
RETURNING *;
