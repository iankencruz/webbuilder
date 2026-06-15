-- +goose Up
CREATE SCHEMA IF NOT EXISTS reference;

CREATE TABLE IF NOT EXISTS reference.block_types (
    code TEXT PRIMARY KEY,
    label TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO reference.block_types 
  (code, label, description) 
VALUES
  ('hero', 'Hero', 'Large banner section with heading, subheading and CTA'),
  ('richtext', 'Rich Text', 'Formatted text content block')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS reference.block_types;
DROP SCHEMA IF EXISTS reference;
