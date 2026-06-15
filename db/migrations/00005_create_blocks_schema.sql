-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS blocks;

CREATE TABLE IF NOT EXISTS blocks.hero (
    id BIGSERIAL PRIMARY KEY,
    heading TEXT NOT NULL,
    subheading TEXT,
    cta_label TEXT,
    cta_url TEXT,
    bg_image_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS trg_blocks_hero_updated_at ON blocks.hero;
CREATE TRIGGER trg_blocks_hero_updated_at
  BEFORE UPDATE ON blocks.hero FOR EACH ROW
  EXECUTE FUNCTION public.set_updated_at();


CREATE TABLE IF NOT EXISTS blocks.richtext (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    format TEXT NOT NULL DEFAULT 'html',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS trg_blocks_richtext_updated_at ON blocks.richtext;
CREATE TRIGGER trg_blocks_richtext_updated_at
  BEFORE UPDATE ON blocks.richtext FOR EACH ROW
  EXECUTE FUNCTION public.set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_blocks_richtext_updated_at ON blocks.richtext;
DROP TABLE IF EXISTS blocks.richtext;

DROP TRIGGER IF EXISTS trg_blocks_hero_updated_at ON blocks.hero;
DROP TABLE IF EXISTS blocks.hero;

DROP SCHEMA IF EXISTS blocks;
-- +goose StatementEnd
