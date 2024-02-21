-- name: ListUsers :many
SELECT * FROM users;

-- name: AddUser :exec
INSERT INTO users (uuid, name, email) VALUES (?, ?, ?);

-- name: ListEvents :many
SELECT
	*,
	0 AS attendees
FROM events
ORDER BY created_at DESC;

-- name: CreateEvent :exec
INSERT INTO events (uuid, description) VALUES (?, ?);

-- name: GetEvent :one
SELECT * FROM events WHERE uuid = ?;

-- name: ListEventAttendees :many
SELECT users.email, users.name, users.uuid, event_attendees.created_at
FROM  event_attendees
INNER JOIN users ON event_attendees.user_uuid = users.uuid
WHERE event_uuid = ?;

-- name: AddAttendee :exec
INSERT INTO event_attendees (event_uuid, user_uuid) VALUES (?, ?);

-- name: MoveAttendees :exec
UPDATE event_attendees SET event_uuid = ? WHERE event_uuid = ?;
