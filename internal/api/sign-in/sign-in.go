package signinapi

import (
	"context"
	"net/http"
	"time"

	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/google"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/api"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/utils"
	usermodel "github.com/nickolasgough/cloud-9-iam/internal/user/model"
	userservice "github.com/nickolasgough/cloud-9-iam/internal/user/service"
)

func SignInWithGoogle(ctx context.Context, gcpClientSecret string, googleService google.Service, authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := googleService.VerifySignInRequest(ctx, r)
		if err != nil {
			statusCode := ierrors.ToHttpStatusCode(err)
			api.WriteError(w, statusCode, err)
			return
		}

		setJWTCookie(w, authService, user)
		http.Redirect(w, r, utils.BuildClientURL("/home"), http.StatusPermanentRedirect)
	}
}

type signInWithPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInWithPassword(ctx context.Context, userService userservice.Service, authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.CorsInterceptor()(w, r, []string{http.MethodPost})
		if r.Method == http.MethodOptions {
			api.WriteError(w, http.StatusOK, nil)
			return
		}

		req := new(signInWithPasswordRequest)
		err := api.UnmarshalRequestBody(r, req)
		if err != nil {
			statusCode := http.StatusBadRequest
			api.WriteError(w, statusCode, err)
			return
		}
		user, err := userService.VerifyUserPassword(req.Email, req.Password)
		if err != nil {
			statusCode := ierrors.ToHttpStatusCode(err)
			api.WriteError(w, statusCode, err)
			return
		}

		setJWTCookie(w, authService, user)
		http.Redirect(w, r, utils.BuildClientURL("/home"), http.StatusPermanentRedirect)
	}
}

func setJWTCookie(w http.ResponseWriter, authService auth.Service, user *usermodel.User) {
	jwtString, err := authService.CreateJWT(user)
	if err != nil {
		statusCode := ierrors.ToHttpStatusCode(err)
		statusText := http.StatusText(statusCode)
		w.WriteHeader(statusCode)
		w.Write([]byte(statusText))
		return
	}
	expiry := time.Now().Add(30 * 24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    jwtString,
		Domain:   "localhost",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteDefaultMode,
		Expires:  expiry,
	})
}
