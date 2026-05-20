-- name: GetUserBySub :one
SELECT * FROM users
WHERE sub = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (sub, provider, email, name, avatar_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    name = $3,
    avatar_url = $4,
    updated_at = NOW()
WHERE sub = $1
RETURNING *;
