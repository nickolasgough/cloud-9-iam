package auth

import usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"

type Service interface {
	CreateJWT(*usermodel.User) (string, error)
	ValidateJWT(string) error
}
