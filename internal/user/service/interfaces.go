package userservice

import usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"

type Service interface {
	CreateUser(*usermodel.User) (*usermodel.User, error)
}
