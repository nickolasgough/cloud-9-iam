package userrepository

import (
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

type Repository interface {
	CreateUser(*usermodel.User, string) (*usermodel.User, error)
	GetUser(string) (*usermodel.User, error)
	GetUserAndPasswordByEmail(string) (*usermodel.User, string, error)
}
