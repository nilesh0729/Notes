-- name: AddTagToNote :exec
INSERT INTO note_tags (note_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromNote :exec
DELETE FROM note_tags
WHERE note_id = $1 AND tag_id = $2;

-- name: GetTagsForNote :many
SELECT t.*
FROM tags t
JOIN note_tags nt ON nt.tag_id = t.id
WHERE nt.note_id = $1;

-- name: GetNotesForTag :many
SELECT n.*
FROM notes n
JOIN note_tags nt ON nt.note_id = n.id
WHERE nt.tag_id = $1;

-- name: DeleteNoteTags :exec
DELETE FROM note_tags
WHERE note_id = $1;
