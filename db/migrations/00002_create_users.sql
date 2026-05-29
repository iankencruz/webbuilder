-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    sub        TEXT NOT NULL UNIQUE,   
    provider   TEXT NOT NULL,         
    email      TEXT NOT NULL,
    first_name TEXT,
    last_name  TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
