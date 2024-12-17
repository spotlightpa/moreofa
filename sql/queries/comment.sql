-- name: CreateComment :one
INSERT INTO
comment (
  name, contact, subject, cc, message, ip, user_agent, referrer, host_page
)
VALUES
(
  ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: ListComments :many
SELECT *
FROM comment
ORDER BY modified_at DESC
LIMIT ? OFFSET ?;
