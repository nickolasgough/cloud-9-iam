package usermodel

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	UsersTable = "users"
)

type User struct {
	ID              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DisplayName     string `json:"displayName"`
	DisplayImageURL string `json:"displayImageURL"`
	Email           string `json:"email"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Deleted time.Time `json:"-"`
}

func Register(db *sql.DB) error {
	_, err := db.Exec(createTable)
	if err != nil {
		return err
	}
	_, err = db.Exec(deletedIndex)
	if err != nil {
		return err
	}
	return nil
}

var createTable = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
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
`, UsersTable)

var deletedIndex = `
CREATE INDEX IF NOT EXISTS users_deleted
ON users (deleted);
`
