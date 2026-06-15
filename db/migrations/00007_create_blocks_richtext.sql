-- +goose Up
CREATE TABLE IF NOT EXISTS blocks.richtext (
  id BIGSERIAL PRIMARY KEY,
  content TEXT NOT NULL DEFAULT '',
  format TEXT NOT NULL DEFAULT 'tiptap_json',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TRIGGER trg_blocks_richtext_updated_at
  BEFORE UPDATE ON blocks.richtext
  FOR EACH ROW
  EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS trg_blocks_richtext_updated_at ON blocks.richtext;
DROP TABLE IF EXISTS blocks.richtext;
