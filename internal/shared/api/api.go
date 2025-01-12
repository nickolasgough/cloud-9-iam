package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nickolasgough/cloud-9-iam/internal/auth"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/constants"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
)

type ApiHandler interface {
	Path() string
	Methods() []string
	Request() interface{}
	Handle(context.Context, interface{}) (interface{}, error)
}

func ServeHTTP(as auth.Service, h ApiHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		CorsInterceptor()(w, r, h.Methods())
		if r.Method == http.MethodOptions {
			WriteError(w, http.StatusOK, nil)
			return
		}

		jwtToken, err := jwtInterceptor(as)(w, r)
		if err != nil {
			WriteError(w, http.StatusUnauthorized, err)
			return
		}
		ctx = context.WithValue(ctx, "session", jwtToken)

		req := h.Request()
		if req != nil {
			err := UnmarshalRequestBody(r, req)
			if err != nil {
				WriteError(w, http.StatusBadRequest, err)
				return
			}
		}
		resp, err := h.Handle(ctx, req)
		if err != nil {
			statusCode := ierrors.ToHttpStatusCode(err)
			WriteError(w, statusCode, err)
			return
		}
		respString := ""
		if resp != nil {
			w.Header().Set("Content-Type", "application/json")
			respString, err = MarshalResponseBody(resp)
			if err != nil {
				w.Header().Set("Content-Type", "text/plain")
				statusCode := http.StatusInternalServerError
				WriteError(w, statusCode, err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(respString))
	}
}

func GetJWTTokenFromContext(ctx context.Context) (*jwt.Token, error) {
	ctxJwt := ctx.Value("session")
	if ctxJwt == nil {
		return nil, ierrors.NewError(ierrors.Unauthorized, errors.New("session is required"))
	}
	jwtToken, ok := ctxJwt.(*jwt.Token)
	if !ok {
		return nil, ierrors.NewError(ierrors.Unauthorized, errors.New("session is malformed"))
	}
	return jwtToken, nil
}

func CorsInterceptor() func(w http.ResponseWriter, r *http.Request, ms []string) {
	return func(w http.ResponseWriter, r *http.Request, ms []string) {
		w.Header().Set("Access-Control-Allow-Origin", constants.CLIENT_BASE_URL)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(ms, ", "))
		if r.Header.Get("Access-Control-Request-Headers") != "" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func jwtInterceptor(authService auth.Service) func(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	return func(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
		jwtCookie, err := r.Cookie("session")
		if err != nil {
			return nil, err
		}
		jwtString := jwtCookie.String()
		jwtToken, err := authService.ValidateJWT(jwtString)
		if err != nil {
			return nil, err
		}
		return jwtToken, nil
	}
}

func UnmarshalRequestBody(req *http.Request, dest interface{}) error {
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(bytes), dest)
}

func MarshalResponseBody(resp interface{}) (string, error) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	statusText := http.StatusText(statusCode)
	w.WriteHeader(statusCode)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s: %s", statusText, err)))
	} else {
		w.Write([]byte(statusText))
	}
}
