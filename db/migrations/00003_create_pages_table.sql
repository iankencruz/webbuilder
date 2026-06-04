-- +goose Up
CREATE TABLE IF NOT EXISTS pages (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    STATUS TEXT NOT NULL DEFAULT 'draft' CHECK (STATUS IN ('draft', 'published', 'archived')),
    seo_title TEXT,
    seo_description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages (slug);
-- +goose StatementBegin
CREATE
OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS
$$
BEGIN
NEW.updated_at = NOW();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_pages_updated_at BEFORE
UPDATE
    ON pages FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose Down
DROP TRIGGER IF EXISTS trg_pages_updated_at ON pages;
DROP FUNCTION IF EXISTS set_updated_at();
DROP INDEX IF EXISTS idx_pages_slug;
DROP TABLE IF EXISTS pages;
