package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nickolasgough/cloud-community-iam/internal/auth"
	"github.com/nickolasgough/cloud-community-iam/internal/google"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/constants"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/ierrors"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/utils"
)

type HandlerFn func(http.ResponseWriter, *http.Request)

func SignInWithGoogle(ctx context.Context, gcpClientSecret string, googleService google.Service, authService auth.Service) HandlerFn {
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
			fmt.Printf("Failed to create JWT: %s\n", err)
			statusCode := ierrors.ToHttpStatusCode(err)
			statusText := http.StatusText(statusCode)
			w.WriteHeader(statusCode)
			w.Write([]byte(statusText))
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    jwtString,
			Domain:   constants.CLIENT_DOMAIN,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		})
		http.Redirect(w, r, utils.BuildClientURL("/home"), http.StatusPermanentRedirect)
	}
}
