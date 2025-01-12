package userrepository

import (
	"database/sql"
)

type repository struct {
	database *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		database: db,
	}
}
