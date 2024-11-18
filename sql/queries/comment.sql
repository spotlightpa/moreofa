-- name: CreateComment :one
INSERT INTO
  comment (
    name,
    contact,
    message,
    ip,
    user_agent,
    referrer,
    host_page
  )
VALUES
  (?, ?, ?, ?, ?, ?, ?) RETURNING *;
