-- +goose Up
CREATE TABLE IF NOT EXISTS blocks_hero (
  id BIGSERIAL PRIMARY KEY,
  heading TEXT NOT NULL,
  subheading TEXT,
  cta_label TEXT,
  cta_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TRIGGER trg_blocks_hero_updated_at
  BEFORE UPDATE ON blocks_hero
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_blocks_hero_updated_at ON blocks_hero;
DROP TABLE IF EXISTS blocks_hero;
