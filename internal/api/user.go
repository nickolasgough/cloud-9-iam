package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
	userservice "github.com/nickolasgough/cloud-9-iam/internal/user/service"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserHandler struct {
	userService userservice.Service
}

func NewCreateUserHandler(us userservice.Service) ApiHandler {
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
		DisplayName: req.Username,
		Password:    req.Password,
	}
	user, err := h.userService.CreateUser(user)
	if err != nil {
		return "", err
	}
	return marshalResponseBody(user)
}
