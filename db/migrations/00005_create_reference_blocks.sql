-- +goose Up
CREATE TABLE IF NOT EXISTS pages_blocks (
    id BIGSERIAL PRIMARY KEY,
    page_id BIGINT NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    block_id BIGINT NOT NULL,
    block_collection TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    hide_block BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_pages_blocks UNIQUE (page_id, block_id, block_collection)
);

CREATE INDEX IF NOT EXISTS idx_pages_blocks_page_order ON pages_blocks (page_id, sort_order);

-- +goose Down
DROP INDEX IF EXISTS idx_pages_blocks_page_order;

DROP TABLE IF EXISTS pages_blocks;
