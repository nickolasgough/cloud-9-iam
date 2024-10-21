package usermodel

import (
	"database/sql"
	"time"
)

type User struct {
	ID              string
	FirstName       string
	LastName        string
	DisplayName     string
	DisplayImageURL string
	Email           string
	Password        string

	Created time.Time
	Updated time.Time
	Deleted time.Time
}

func Register(db *sql.DB) error {
	_, err := db.Exec(createTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(createDeletedIndex)
	if err != nil {
		return err
	}
	return nil
}

var createTable = `
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(250) PRIMARY KEY,
	firstName VARCHAR(250),
	lastName VARCHAR(250),
	displayName VARCHAR(250),
	displayImageURL VARCHAR(250),
	email VARCHAR(250) UNIQUE NOT NULL,
	password VARCHAR(250) NOT NULL,
	created TIMESTAMPTZ NOT NULL,
	updated TIMESTAMPTZ NOT NULL,
	deleted TIMESTAMPTZ
);
`

var createDeletedIndex = `
CREATE INDEX IF NOT EXISTS users_deleted
ON users (deleted);
`
