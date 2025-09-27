-- name: CreateTags :one
INSERT INTO Tags (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: GetTag :one
SELECT * FROM tags
WHERE tag_id = $1 
LIMIT 1;

-- name: ListTags :many
SELECT * FROM tags
WHERE tag_id > $1
ORDER BY tag_id
LIMIT $2;

-- name: UpdateTag :one
UPDATE Tags
SET name = $2
WHERE tag_id = $1
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM Tags
WHERE tag_id = $1;