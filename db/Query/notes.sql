-- name: CreateNotes :one
INSERT INTO notes (
  user_id,
  title,
  content, 
  pinned,
  archived
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetNoteById :one
SELECT * FROM notes
WHERE id = $1;

-- name: ListUserNotes :many
SELECT * FROM notes
WHERE user_id = $1 AND archived = FALSE 
ORDER BY created_at DESC;


-- name: UpdateNotes :one
UPDATE notes
  set title = $2,
      content = $3,
      pinned = $4,
      archived = $5,
      updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteNotes :exec
DELETE FROM notes
WHERE id = $1;