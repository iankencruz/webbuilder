-- name: CreatePage :one
INSERT INTO
    pages (
        title,
        slug,
        STATUS,
        seo_title,
        seo_description
    )
VALUES
    (
        @title,
        @slug,
        @status,
        @seo_title,
        @seo_description
    )
RETURNING
    *;
-- name: GetPageBySlug :one
SELECT
    *
FROM
    pages
WHERE
    slug = @slug;
-- name: GetPageByID :one
SELECT
    *
FROM
    pages
WHERE
    id = @id;
-- name: ListPages :many
SELECT
    *
FROM
    pages
ORDER BY
    created_at DESC;
-- name: UpdatePage :one
UPDATE
    pages
SET
    title = @title,
    slug = @slug,
    STATUS = @status,
    seo_title = @seo_title,
    seo_description = @seo_description
WHERE
    slug = @slug
RETURNING
    *;
-- name: DeleteBySlug :exec
DELETE FROM
    pages
WHERE
    slug = @slug;
-- name: DeleteByID :exec
DELETE FROM
    pages
WHERE
    id = @id;
