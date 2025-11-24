-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    NULL
);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: ResetRefreshTokens :exec
DELETE FROM refresh_tokens;