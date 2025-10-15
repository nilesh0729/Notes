-- name: AddTagToNote :one
INSERT INTO note_tags (note_id, tag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING note_id, tag_id;



-- name: GetTagsForNote :many
SELECT t.tag_id, t.name
FROM tags t
INNER JOIN note_tags nt ON t.tag_id = nt.tag_id
WHERE nt.note_id = $1;


-- name: GetNotesForTag :many
SELECT n.note_id, n.title, n.owner, n.content, n.pinned, n.archived, n.created_at, n.updated_at
FROM notes n
INNER JOIN note_tags nt ON n.note_id = nt.note_id
WHERE nt.tag_id = $1;



-- name: RemoveTagFromNote :exec
DELETE FROM note_tags
WHERE note_id = $1 AND tag_id = $2;
