-- name: CreateNote :one
INSERT INTO notes (
  title,
  content
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetNoteById :one
SELECT * FROM notes
WHERE note_id = $1 
LIMIT 1;

-- name: ListNotes :many
SELECT * FROM notes
WHERE note_id > $1
ORDER BY note_id 
LIMIT $2;

-- name: UpdateNote :one
UPDATE notes
  set title = $2,
  content = $3
WHERE note_id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE note_id = $1;