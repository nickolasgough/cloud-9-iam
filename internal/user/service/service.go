package userservice

import (
	"errors"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
	userrepository "github.com/nickolasgough/cloud-9-iam/internal/user/repository"
)

type service struct {
	repository userrepository.Repository
}

func NewService(r userrepository.Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) CreateUser(user *usermodel.User, password string) (*usermodel.User, error) {
	return s.repository.CreateUser(user, password)
}

func (s *service) GetUser(userID string) (*usermodel.User, error) {
	return s.repository.GetUser(userID)
}

func (s *service) VerifyUserPassword(email string, submittedPassword string) (*usermodel.User, error) {
	user, password, err := s.repository.GetUserAndPasswordByEmail(email)
	if err != nil {
		return nil, err
	}
	if submittedPassword != password {
		return nil, ierrors.NewError(ierrors.Forbidden, errors.New("invalid user credentials"))
	}
	return user, nil
}
