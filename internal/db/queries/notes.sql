-- name: CreateNote :one
INSERT INTO notes (
  owner,
  title,
  content
) VALUES (
  $1, $2 ,$3
)
RETURNING *;

-- name: GetNoteById :one
SELECT * FROM notes
WHERE note_id = $1
LIMIT 1;

-- name: ListNotes :many
SELECT * FROM notes
WHERE note_id > $1 AND owner = $3
ORDER BY note_id 
LIMIT $2;

-- name: SearchNotes :many
SELECT * FROM notes
WHERE (title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%') AND owner = $4
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateNote :one
UPDATE notes
  set title = $2,
  content = $3
WHERE note_id = $1
RETURNING *;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE note_id = $1;

-- name: DeleteNoteTagsByNoteId :exec
DELETE FROM note_tags
WHERE note_id = $1;