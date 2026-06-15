-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pages_blocks (
    id BIGSERIAL PRIMARY KEY,
    page_id BIGINT NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    block_id BIGINT NOT NULL,
    block_collection TEXT NOT NULL REFERENCES reference.block_types(code),
    sort_order INT NOT NULL DEFAULT 0,
    hide_block BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_pages_blocks UNIQUE (page_id, block_id, block_collection)
);

CREATE INDEX IF NOT EXISTS idx_pages_blocks_page_order ON pages_blocks (page_id, sort_order);

DROP TRIGGER IF EXISTS trg_pages_blocks_updated_at ON pages_blocks;
CREATE TRIGGER trg_pages_blocks_updated_at
  BEFORE UPDATE ON pages_blocks FOR EACH ROW
  EXECUTE FUNCTION public.set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_pages_blocks_updated_at ON pages_blocks;
DROP INDEX IF EXISTS idx_pages_blocks_page_order;
DROP TABLE IF EXISTS pages_blocks;
-- +goose StatementEnd
