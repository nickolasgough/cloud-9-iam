package userservice

import usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"

type Service interface {
	CreateUser(*usermodel.User, string) (*usermodel.User, error)
	GetUser(string) (*usermodel.User, error)
	VerifyUserPassword(string, string) (*usermodel.User, error)
}
