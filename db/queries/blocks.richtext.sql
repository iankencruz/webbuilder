-- name: CreateRichTextBlock :one
INSERT INTO
    blocks.richtext (content, format)
VALUES
    (@content, @format)
RETURNING
    *;
-- name: GetRichTextBlock :one
SELECT
    *
FROM
    blocks.richtext
WHERE
    id = @id;
-- name: UpdateRichTextBlock :one
UPDATE
    blocks.richtext
SET
    content = @content,
    format = @format
WHERE
    id = @id
RETURNING
    *;
-- name: DeleteRichTextBlock :exec
DELETE FROM
    blocks.richtext
WHERE
    id = @id;
