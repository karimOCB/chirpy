-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id) 
VALUES ( 
    $1,
    NOW(),
    NOW(),
    $2,
    NULL,
    $3
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :execresult
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;