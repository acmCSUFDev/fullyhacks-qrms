-- name: ListUsers :many
SELECT * FROM users;

-- name: AddUser :exec
INSERT INTO users (code, name, email) VALUES (?, ?, ?);

-- name: GetUser :one
SELECT name FROM users WHERE email = ?;

-- name: GetUserCode :one
SELECT code FROM users WHERE email = ?;

-- name: ListEvents :many
SELECT
	*,
	(SELECT COUNT(*) FROM event_attendees WHERE event_uuid = events.uuid) AS attendees
FROM events
ORDER BY created_at DESC;

-- name: CreateEvent :exec
INSERT INTO events (uuid, description) VALUES (?, ?);

-- name: GetEvent :one
SELECT
	*,
	(SELECT COUNT(*) FROM event_attendees WHERE event_uuid = events.uuid) AS attendees
FROM events
WHERE uuid = ?;

-- name: ListEventAttendees :many
SELECT users.email, users.name, event_attendees.created_at
FROM  event_attendees
INNER JOIN users ON event_attendees.user_code = users.code
WHERE event_uuid = ?;

-- name: AddAttendee :exec
INSERT INTO event_attendees (event_uuid, user_code) VALUES (?, ?);

-- name: MoveAttendees :exec
UPDATE event_attendees SET event_uuid = ? WHERE event_uuid = ?;

-- name: AddAuthToken :one
INSERT INTO auth_tokens (token, parent_token) VALUES (?, ?) RETURNING created_at;

-- name: CheckAuthToken :one
SELECT parent_token, created_at FROM auth_tokens WHERE token = ?;
