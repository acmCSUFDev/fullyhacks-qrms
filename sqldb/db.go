// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqldb

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.addAttendeeStmt, err = db.PrepareContext(ctx, addAttendee); err != nil {
		return nil, fmt.Errorf("error preparing query AddAttendee: %w", err)
	}
	if q.addUserStmt, err = db.PrepareContext(ctx, addUser); err != nil {
		return nil, fmt.Errorf("error preparing query AddUser: %w", err)
	}
	if q.createEventStmt, err = db.PrepareContext(ctx, createEvent); err != nil {
		return nil, fmt.Errorf("error preparing query CreateEvent: %w", err)
	}
	if q.getEventStmt, err = db.PrepareContext(ctx, getEvent); err != nil {
		return nil, fmt.Errorf("error preparing query GetEvent: %w", err)
	}
	if q.listEventAttendeesStmt, err = db.PrepareContext(ctx, listEventAttendees); err != nil {
		return nil, fmt.Errorf("error preparing query ListEventAttendees: %w", err)
	}
	if q.listEventsStmt, err = db.PrepareContext(ctx, listEvents); err != nil {
		return nil, fmt.Errorf("error preparing query ListEvents: %w", err)
	}
	if q.listUsersStmt, err = db.PrepareContext(ctx, listUsers); err != nil {
		return nil, fmt.Errorf("error preparing query ListUsers: %w", err)
	}
	if q.moveAttendeesStmt, err = db.PrepareContext(ctx, moveAttendees); err != nil {
		return nil, fmt.Errorf("error preparing query MoveAttendees: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.addAttendeeStmt != nil {
		if cerr := q.addAttendeeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addAttendeeStmt: %w", cerr)
		}
	}
	if q.addUserStmt != nil {
		if cerr := q.addUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing addUserStmt: %w", cerr)
		}
	}
	if q.createEventStmt != nil {
		if cerr := q.createEventStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createEventStmt: %w", cerr)
		}
	}
	if q.getEventStmt != nil {
		if cerr := q.getEventStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEventStmt: %w", cerr)
		}
	}
	if q.listEventAttendeesStmt != nil {
		if cerr := q.listEventAttendeesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listEventAttendeesStmt: %w", cerr)
		}
	}
	if q.listEventsStmt != nil {
		if cerr := q.listEventsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listEventsStmt: %w", cerr)
		}
	}
	if q.listUsersStmt != nil {
		if cerr := q.listUsersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUsersStmt: %w", cerr)
		}
	}
	if q.moveAttendeesStmt != nil {
		if cerr := q.moveAttendeesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing moveAttendeesStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                     DBTX
	tx                     *sql.Tx
	addAttendeeStmt        *sql.Stmt
	addUserStmt            *sql.Stmt
	createEventStmt        *sql.Stmt
	getEventStmt           *sql.Stmt
	listEventAttendeesStmt *sql.Stmt
	listEventsStmt         *sql.Stmt
	listUsersStmt          *sql.Stmt
	moveAttendeesStmt      *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                     tx,
		tx:                     tx,
		addAttendeeStmt:        q.addAttendeeStmt,
		addUserStmt:            q.addUserStmt,
		createEventStmt:        q.createEventStmt,
		getEventStmt:           q.getEventStmt,
		listEventAttendeesStmt: q.listEventAttendeesStmt,
		listEventsStmt:         q.listEventsStmt,
		listUsersStmt:          q.listUsersStmt,
		moveAttendeesStmt:      q.moveAttendeesStmt,
	}
}
