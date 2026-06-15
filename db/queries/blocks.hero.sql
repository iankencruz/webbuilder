-- name: CreateHeroBlock :one
INSERT INTO
    blocks.hero (
        heading,
        subheading,
        cta_label,
        cta_url
    )
VALUES
    (
        @heading,
        @subheading,
        @cta_label,
        @cta_url
    )
RETURNING
    *;
-- name: GetHeroBlock :one
SELECT
    *
FROM
    blocks.hero
WHERE
    id = @id;
-- name: UpdateHeroBlock :one
UPDATE
    blocks.hero
SET
    heading = @heading,
    subheading = @subheading,
    cta_label = @cta_label,
    cta_url = @cta_url
WHERE
    id = @id
RETURNING
    *;
-- name: DeleteHeroBlock :exec
DELETE FROM
    blocks.hero
WHERE
    id = @id;
