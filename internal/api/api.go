package api

import (
	"context"
	"net/http"

	"github.com/nickolasgough/cloud-community-iam/internal/google"
	"github.com/nickolasgough/cloud-community-iam/internal/shared/ierrors"
)

type HandlerFn func(http.ResponseWriter, *http.Request)

func SignInWithGoogle(ctx context.Context, googleService google.Service) HandlerFn {
	return func(w http.ResponseWriter, r *http.Request) {
		err := googleService.VerifySignInRequest(ctx, r)
		if err != nil {
			statusCode := ierrors.ToHttpStatusCode(err)
			// message := http.StatusText(statusCode)
			w.WriteHeader(statusCode)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "https://localhost:4200/home", http.StatusPermanentRedirect)
	}
}
