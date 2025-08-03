-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;


-- name: GetUsers :one
SELECT * FROM users
WHERE id = $1 
LIMIT 1;


-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;


-- name: UpdateUsers :one
UPDATE users
  set password_hash = $2,
  email = $3
WHERE id = $1
RETURNING *;


-- name: DeleteUsers :exec
DELETE FROM users
WHERE id = $1;