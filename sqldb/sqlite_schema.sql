CREATE TABLE users (
	uuid TEXT PRIMARY KEY, -- unique identifier
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE);

CREATE TABLE events (
	uuid TEXT PRIMARY KEY, -- unique identifier
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	description TEXT NOT NULL);

CREATE TABLE event_attendees (
	event_uuid TEXT NOT NULL,
	user_uuid TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (event_uuid, user_uuid),
	FOREIGN KEY (event_uuid) REFERENCES events(uuid),
	FOREIGN KEY (user_uuid) REFERENCES users(uuid));

--------------------------------- NEW VERSION ---------------------------------

CREATE TABLE auth_tokens (
	token TEXT PRIMARY KEY, -- unique identifier
	parent_token TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);
