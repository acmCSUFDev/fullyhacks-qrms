CREATE TABLE users (
	email TEXT PRIMARY KEY,
	code TEXT NOT NULL UNIQUE, -- unique identifier
	name TEXT NOT NULL);

CREATE TABLE events (
	uuid TEXT PRIMARY KEY, -- unique identifier
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	description TEXT NOT NULL);

CREATE TABLE event_attendees (
	event_uuid TEXT NOT NULL,
	user_code TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (event_uuid, user_code),
	FOREIGN KEY (event_uuid) REFERENCES events(uuid),
	FOREIGN KEY (user_code) REFERENCES users(code));

CREATE TABLE auth_tokens (
	token TEXT PRIMARY KEY, -- unique identifier
	parent_token TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
