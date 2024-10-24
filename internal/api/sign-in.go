package api

import (
	"context"
	"net/http"
	"time"

	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/google"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/utils"
)

func SignInWithGoogle(ctx context.Context, gcpClientSecret string, googleService google.Service, authService auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := googleService.VerifySignInRequest(ctx, r)
		if err != nil {
			statusCode := ierrors.ToHttpStatusCode(err)
			statusText := http.StatusText(statusCode)
			w.WriteHeader(statusCode)
			w.Write([]byte(statusText))
			return
		}

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
		http.Redirect(w, r, utils.BuildClientURL("/home"), http.StatusPermanentRedirect)
	}
}
