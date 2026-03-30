-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), -- automatically generate UUID
    NOW(),             -- created_at timestamp
    NOW(),             -- updated_at timestamp
    $1,                -- email parameter
    $2                 -- hashed_password parameter
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET hashed_password = $1,
    email = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeUser :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
