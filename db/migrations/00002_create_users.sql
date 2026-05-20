-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    sub        TEXT NOT NULL UNIQUE,   -- stable OIDC subject ID
    provider   TEXT NOT NULL,          -- "rauthy", "google" etc
    email      TEXT NOT NULL,
    name       TEXT,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
