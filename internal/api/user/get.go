package userapi

import (
	"context"
	"net/http"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/api"
	userservice "github.com/nickolasgough/cloud-9-iam/internal/user/service"
)

type getUserHandler struct {
	userService userservice.Service
}

func NewGetUserHandler(us userservice.Service) api.ApiHandler {
	return &getUserHandler{
		userService: us,
	}
}

func (h *getUserHandler) Path() string {
	return "Get /user"
}

func (h *getUserHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getUserHandler) Request() interface{} {
	return nil
}

func (h *getUserHandler) Handle(ctx context.Context, r interface{}) (interface{}, error) {
	jwtToken, err := api.GetJWTTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		return "", err
	}
	return user, nil
}
