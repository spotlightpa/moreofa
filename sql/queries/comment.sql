-- name: CreateComment :one
INSERT INTO
comment (subject, name, contact, message, ip, user_agent, referrer, host_page)
VALUES
(
  ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING * ;
