-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,                 -- token
    NOW(),              -- created_at timestamp
    NOW(),              -- updated_at timestamp
    $2,                 -- user_id
    $3                  -- expires_at
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users AS u
LEFT JOIN refresh_tokens AS rt
    ON u.id = rt.user_id
WHERE rt.expires_at > NOW()
AND rt.revoked_at IS NULL
AND rt.token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;

