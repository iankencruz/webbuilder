-- name: CreateRichTextBlock :one
INSERT INTO
    blocks_richtext (content, format)
VALUES
    (@content, @format)
RETURNING
    *;
-- name: GetRichTextBlock :one
SELECT
    *
FROM
    blocks_richtext
WHERE
    id = @id;
-- name: UpdateRichTextBlock :one
UPDATE
    blocks_richtext
SET
    content = @content,
    format = @format
WHERE
    id = @id
RETURNING
    *;
-- name: DeleteRichTextBlock :exec
DELETE FROM
    blocks_richtext
WHERE
    id = @id;
