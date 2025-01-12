package auth

import (
	"github.com/golang-jwt/jwt/v5"

	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
)

type Service interface {
	CreateJWT(*usermodel.User) (string, error)
	ValidateJWT(string) (*jwt.Token, error)
}
