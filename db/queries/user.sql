-- name: GetUserBySub :one
SELECT * FROM users
WHERE sub = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (sub, provider, email, first_name, last_name, avatar_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $2,
    first_name = $3,
    last_name = $4,
    avatar_url = $5,
    updated_at = NOW()
WHERE sub = $1
RETURNING *;
