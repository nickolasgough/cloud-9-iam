package userrepository

import (
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

type Repository interface {
	CreateUser(user *usermodel.User) (*usermodel.User, error)
}
