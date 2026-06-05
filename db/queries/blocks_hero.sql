-- name: CreateHeroBlock :one
INSERT INTO
    blocks_hero (
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
    blocks_hero
WHERE
    id = @id;
-- name: UpdateHeroBlock :one
UPDATE
    blocks_hero
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
    blocks_hero
WHERE
    id = @id;
