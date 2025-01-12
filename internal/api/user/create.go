package userapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/api"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
	userservice "github.com/nickolasgough/cloud-9-iam/internal/user/service"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserHandler struct {
	userService userservice.Service
}

func NewCreateUserHandler(us userservice.Service) api.ApiHandler {
	return &createUserHandler{
		userService: us,
	}
}

func (h *createUserHandler) Path() string {
	return "/user"
}

func (h *createUserHandler) Methods() []string {
	return []string{http.MethodPost}
}

func (h *createUserHandler) Request() interface{} {
	return new(createUserRequest)
}

func (h *createUserHandler) Handle(ctx context.Context, r interface{}) (interface{}, error) {
	req, ok := r.(*createUserRequest)
	if !ok {
		return nil, ierrors.NewError(ierrors.InvalidArgument, errors.New("received wrong request type"))
	}

	user := &usermodel.User{
		Email: req.Email,
	}
	user, err := h.userService.CreateUser(user, req.Password)
	if err != nil {
		return "", err
	}
	return user, nil
}
