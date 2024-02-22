// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: sqlite_queries.sql

package sqldb

import (
	"context"
	"time"
)

const addAttendee = `-- name: AddAttendee :exec
INSERT INTO event_attendees (event_uuid, user_uuid) VALUES (?, ?)
`

type AddAttendeeParams struct {
	EventUUID string
	UserUUID  string
}

func (q *Queries) AddAttendee(ctx context.Context, arg AddAttendeeParams) error {
	_, err := q.exec(ctx, q.addAttendeeStmt, addAttendee, arg.EventUUID, arg.UserUUID)
	return err
}

const addAuthToken = `-- name: AddAuthToken :one
INSERT INTO auth_tokens (token, parent_token) VALUES (?, ?) RETURNING created_at
`

type AddAuthTokenParams struct {
	Token       string
	ParentToken string
}

func (q *Queries) AddAuthToken(ctx context.Context, arg AddAuthTokenParams) (time.Time, error) {
	row := q.queryRow(ctx, q.addAuthTokenStmt, addAuthToken, arg.Token, arg.ParentToken)
	var created_at time.Time
	err := row.Scan(&created_at)
	return created_at, err
}

const addUser = `-- name: AddUser :exec
INSERT INTO users (uuid, name, email) VALUES (?, ?, ?)
`

type AddUserParams struct {
	UUID  string
	Name  string
	Email string
}

func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) error {
	_, err := q.exec(ctx, q.addUserStmt, addUser, arg.UUID, arg.Name, arg.Email)
	return err
}

const checkAuthToken = `-- name: CheckAuthToken :one
SELECT parent_token, created_at FROM auth_tokens WHERE token = ?
`

type CheckAuthTokenRow struct {
	ParentToken string
	CreatedAt   time.Time
}

func (q *Queries) CheckAuthToken(ctx context.Context, token string) (CheckAuthTokenRow, error) {
	row := q.queryRow(ctx, q.checkAuthTokenStmt, checkAuthToken, token)
	var i CheckAuthTokenRow
	err := row.Scan(&i.ParentToken, &i.CreatedAt)
	return i, err
}

const createEvent = `-- name: CreateEvent :exec
INSERT INTO events (uuid, description) VALUES (?, ?)
`

type CreateEventParams struct {
	UUID        string
	Description string
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) error {
	_, err := q.exec(ctx, q.createEventStmt, createEvent, arg.UUID, arg.Description)
	return err
}

const getEvent = `-- name: GetEvent :one
SELECT uuid, created_at, description FROM events WHERE uuid = ?
`

func (q *Queries) GetEvent(ctx context.Context, uuid string) (Event, error) {
	row := q.queryRow(ctx, q.getEventStmt, getEvent, uuid)
	var i Event
	err := row.Scan(&i.UUID, &i.CreatedAt, &i.Description)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT uuid, name, email FROM users WHERE uuid = ?
`

func (q *Queries) GetUser(ctx context.Context, uuid string) (User, error) {
	row := q.queryRow(ctx, q.getUserStmt, getUser, uuid)
	var i User
	err := row.Scan(&i.UUID, &i.Name, &i.Email)
	return i, err
}

const listEventAttendees = `-- name: ListEventAttendees :many
SELECT users.email, users.name, users.uuid, event_attendees.created_at
FROM  event_attendees
INNER JOIN users ON event_attendees.user_uuid = users.uuid
WHERE event_uuid = ?
`

type ListEventAttendeesRow struct {
	Email     string
	Name      string
	UUID      string
	CreatedAt time.Time
}

func (q *Queries) ListEventAttendees(ctx context.Context, eventUuid string) ([]ListEventAttendeesRow, error) {
	rows, err := q.query(ctx, q.listEventAttendeesStmt, listEventAttendees, eventUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListEventAttendeesRow
	for rows.Next() {
		var i ListEventAttendeesRow
		if err := rows.Scan(
			&i.Email,
			&i.Name,
			&i.UUID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEvents = `-- name: ListEvents :many
SELECT
	uuid, created_at, description,
	(SELECT COUNT(*) FROM event_attendees WHERE event_uuid = events.uuid) AS attendees
FROM events
ORDER BY created_at DESC
`

type ListEventsRow struct {
	UUID        string
	CreatedAt   time.Time
	Description string
	Attendees   int64
}

func (q *Queries) ListEvents(ctx context.Context) ([]ListEventsRow, error) {
	rows, err := q.query(ctx, q.listEventsStmt, listEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListEventsRow
	for rows.Next() {
		var i ListEventsRow
		if err := rows.Scan(
			&i.UUID,
			&i.CreatedAt,
			&i.Description,
			&i.Attendees,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT uuid, name, email FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.listUsersStmt, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.UUID, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const moveAttendees = `-- name: MoveAttendees :exec
UPDATE event_attendees SET event_uuid = ? WHERE event_uuid = ?
`

type MoveAttendeesParams struct {
	EventUUID   string
	EventUUID_2 string
}

func (q *Queries) MoveAttendees(ctx context.Context, arg MoveAttendeesParams) error {
	_, err := q.exec(ctx, q.moveAttendeesStmt, moveAttendees, arg.EventUUID, arg.EventUUID_2)
	return err
}
