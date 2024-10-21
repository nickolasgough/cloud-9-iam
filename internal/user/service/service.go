package userservice

import (
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

func (s *service) CreateUser(user *usermodel.User) (*usermodel.User, error) {
	return s.repository.CreateUser(user)
}
