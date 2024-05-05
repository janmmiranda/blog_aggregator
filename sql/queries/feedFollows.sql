-- name: CreateFeedFollow :one
INSERT INTO feedFollows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFeedFollow :exec
DELETE FROM feedFollows WHERE id = $1;

-- name: GetFeedFollows :many
SELECT * FROM feedFollows WHERE user_id = $1;