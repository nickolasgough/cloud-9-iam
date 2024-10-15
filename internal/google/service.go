package google

import (
	"context"
	"errors"
	"net/http"

	"github.com/nickolasgough/cloud-community-iam/internal/shared/ierrors"
	usermodel "github.com/nickolasgough/cloud-community-iam/internal/user"
)

const (
	CSRF_TOKEN_NAME = "g_csrf_token"
)

type Service interface {
	VerifySignInRequest(ctx context.Context, req *http.Request) (*usermodel.User, error)
}

type service struct {
	gcpAuthClientID string
	googleIDVerfier GoogleIDVerifier
}

func NewService(gaci string, giv GoogleIDVerifier) Service {
	return &service{
		gcpAuthClientID: gaci,
		googleIDVerfier: giv,
	}
}

func (s *service) VerifySignInRequest(ctx context.Context, req *http.Request) (*usermodel.User, error) {
	err := s.verifyCSRFTokens(req)
	if err != nil {
		return nil, err
	}

	err = req.ParseForm()
	if err != nil {
		return nil, ierrors.NewError(ierrors.BadRequest, err)
	}
	googleJWT := req.Form.Get("credential")
	if googleJWT == "" {
		return nil, ierrors.NewError(ierrors.Unauthorized, err)
	}
	payload, err := s.googleIDVerfier.Validate(ctx, googleJWT, s.gcpAuthClientID)
	if err != nil {
		return nil, ierrors.NewError(ierrors.Forbidden, err)
	}

	user, err := extractUserFromJWTClaims(payload.Claims)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) verifyCSRFTokens(req *http.Request) error {
	csrfCookie, err := req.Cookie(CSRF_TOKEN_NAME)
	if err != nil {
		return err
	}
	if csrfCookie == nil || csrfCookie.Value == "" {
		err = errors.New("CSRF token missing in cookie")
		return ierrors.NewError(ierrors.BadRequest, err)
	}

	err = req.ParseForm()
	if err != nil {
		return err
	}
	csrfBody := req.Form.Get(CSRF_TOKEN_NAME)
	if csrfBody == "" {
		err = errors.New("CSRF token missing in post body")
		return ierrors.NewError(ierrors.BadRequest, err)
	}

	if csrfCookie.Value != csrfBody {
		err = errors.New("CSRF tokens don't match")
		return ierrors.NewError(ierrors.BadRequest, err)
	}
	return nil
}

func extractUserFromJWTClaims(jwtClaims map[string]interface{}) (*usermodel.User, error) {
	firstName := jwtClaims["given_name"].(string)
	lastName := jwtClaims["family_name"].(string)
	displayName := jwtClaims["name"].(string)
	displayImageURL := jwtClaims["picture"].(string)
	email := jwtClaims["email"].(string)
	if displayName == "" || email == "" {
		return nil, ierrors.NewError(ierrors.InvalidArgument, errors.New("display name and email are required"))
	}
	return &usermodel.User{
		FirstName:       firstName,
		LastName:        lastName,
		DisplayName:     displayName,
		DisplayImageURL: displayImageURL,
		Email:           email,
	}, nil
}
