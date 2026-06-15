-- +goose Up
-- +goose StatementBegin
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

DROP TRIGGER IF EXISTS trg_pages_updated_at ON pages;
CREATE TRIGGER trg_pages_updated_at 
  BEFORE UPDATE ON pages FOR EACH ROW
  EXECUTE FUNCTION public.set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_pages_updated_at ON pages;
DROP INDEX IF EXISTS idx_pages_slug;
DROP TABLE IF EXISTS pages;
-- +goose StatementEnd
