-- +goose Up
CREATE SCHEMA IF NOT EXISTS blocks;
CREATE SCHEMA IF NOT EXISTS reference;



CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP INDEX IF EXISTS sessions_expiry_idx;
DROP TABLE IF EXISTS sessions;


DROP SCHEMA IF EXISTS reference;
DROP SCHEMA IF EXISTS blocks;
