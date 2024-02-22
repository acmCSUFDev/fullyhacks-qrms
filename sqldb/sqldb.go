package sqldb

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	_ "embed"

	"github.com/google/uuid"
	"libdb.so/lazymigrate"

	_ "modernc.org/sqlite"
)

//go:generate sqlc generate

//go:embed sqlite_schema.sql
var schema string

const pragma = `
PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA strict = ON;
`

// Database provides methods for interacting with the database.
// For now, it just wraps around sqlc's Queries because I'm lazy.
type Database struct {
	*Queries
	db *sql.DB
}

// Open creates a new database at the given path.
func Open(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return newDatabase(db)
}

// NewInMemory creates a new in-memory database.
func NewInMemory() (*Database, error) {
	db, _ := sql.Open("sqlite", ":memory:")
	return newDatabase(db)
}

func newDatabase(db *sql.DB) (*Database, error) {
	if _, err := db.Exec(pragma); err != nil {
		return nil, err
	}

	schema := lazymigrate.NewSchema(schema)
	if err := schema.Migrate(context.Background(), db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &Database{
		Queries: New(db),
		db:      db,
	}, nil
}

// Close closes the database.
func (db *Database) Close() error {
	return db.db.Close()
}

// Tx scopes f to a transaction.
func (db *Database) Tx(f func(*Queries) error) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := f(db.Queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit()
}

// GenerateUUID generates a new UUID.
func GenerateUUID() string { return uuid.New().String() }

// NewUser creates a new user.
func NewUser(name, email string) User {
	var randBytes [4]byte
	_, err := rand.Read(randBytes[:])
	if err != nil {
		panic(err)
	}
	randBits := base64.RawURLEncoding.EncodeToString(randBytes[:])

	userHash := sha256.Sum256([]byte(email))
	userBits := base64.RawURLEncoding.EncodeToString(userHash[:])[:6]
	code := "fullyhacks:" + randBits + userBits

	return User{
		Email: email,
		Code:  code,
		Name:  name,
	}
}
