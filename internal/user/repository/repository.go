package userrepository

import (
	"database/sql"

	"github.com/google/uuid"

	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

type repository struct {
	database *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		database: db,
	}
}

func (r *repository) CreateUser(user *usermodel.User) (*usermodel.User, error) {
	user.ID = uuid.NewString()
	return user, nil
}
