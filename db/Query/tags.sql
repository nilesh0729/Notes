-- name: CreateTags :one
INSERT INTO tags (
  name,
  user_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUserTag :one
SELECT * FROM tags
WHERE user_id = $1 AND name = $2
LIMIT 1;

-- name: ListUserTags :many
SELECT * FROM tags
WHERE user_id = $1
ORDER BY name;

-- name: RenameTag :one
UPDATE tags
  set name = $3
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteTags :exec
DELETE FROM tags
WHERE id = $1 AND user_id = $2;