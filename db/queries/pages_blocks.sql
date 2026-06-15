-- name: AddBlockToPage :one
INSERT INTO
    pages_blocks (
        page_id,
        block_id,
        block_collection,
        sort_order
    )
VALUES
    (
        @page_id,
        @block_id,
        @block_collection,
        @sort_order
    )
RETURNING
    *;
-- name: GetPageBlocks :many
SELECT
    *
FROM
    pages_blocks
WHERE
    page_id = @page_id
ORDER BY
    sort_order ASC;
-- name: UpdatePageBlock :one
UPDATE
    pages_blocks
SET
    hide_block = @hide_block,
    sort_order = @sort_order
WHERE
    id = @id
RETURNING
    *;
-- name: DeletePageBlock :exec
DELETE FROM
    pages_blocks
WHERE
    id = @id;
-- name: ReorderPageBlocks :exec
UPDATE
    pages_blocks
SET
    sort_order = @sort_order
WHERE
    id = @id;


-- name: GetPageBlock :one
SELECT
    *
FROM
    pages_blocks
WHERE
    id = @id;
