package google

import (
	"context"
	"errors"
	"net/http"

	"github.com/nickolasgough/cloud-community-iam/internal/shared/ierrors"
)

const (
	CSRF_TOKEN_NAME = "g_csrf_token"
)

type Service interface {
	VerifySignInRequest(ctx context.Context, req *http.Request) error
}

type service struct {
	gcpAuthClientID string
	googleIDVerfier GoogleIDVerifier
}

func NewService(gac string, giv GoogleIDVerifier) Service {
	return &service{
		gcpAuthClientID: gac,
		googleIDVerfier: giv,
	}
}

func (s *service) VerifySignInRequest(ctx context.Context, req *http.Request) error {
	err := s.verifyCSRFTokens(req)
	if err != nil {
		return err
	}

	err = req.ParseForm()
	if err != nil {
		return ierrors.NewError(ierrors.BadRequest, err)
	}
	googleJWT := req.Form.Get("credential")
	if googleJWT == "" {
		return ierrors.NewError(ierrors.Unauthorized, err)
	}
	_, err = s.googleIDVerfier.Validate(ctx, googleJWT, s.gcpAuthClientID)
	if err != nil {
		return ierrors.NewError(ierrors.Forbidden, err)
	}
	return nil
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
