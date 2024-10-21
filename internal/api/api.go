package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/nickolasgough/cloud-9-iam/internal/shared/constants"
	"github.com/nickolasgough/cloud-9-iam/internal/shared/ierrors"
)

type ApiHandler interface {
	Path() string
	Methods() []string
	Request() interface{}
	Handle(context.Context, interface{}) (interface{}, error)
}

type ApiEndpoint struct {
	ApiHandler
}

func NewApiEndpoint(ah ApiHandler) *ApiEndpoint {
	return &ApiEndpoint{
		ApiHandler: ah,
	}
}

func (a *ApiEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	corsInterceptor(w, r, a.ApiHandler)

	req := a.Request()
	if req != nil {
		err := unmarshalRequestBody(r, req)
		if err != nil {
			statusCode := http.StatusBadRequest
			statusText := http.StatusText(statusCode)
			w.WriteHeader(statusCode)
			w.Write([]byte(fmt.Sprintf("%s: %s", statusText, err)))
			return
		}
	}
	resp, err := a.Handle(r.Context(), req)
	if err != nil {
		statusCode := ierrors.ToHttpStatusCode(err)
		statusText := http.StatusText(statusCode)
		w.WriteHeader(statusCode)
		w.Write([]byte(fmt.Sprintf("%s: %s", statusText, err)))
		return
	}
	respString := ""
	if resp != nil {
		respString, err = marshalResponseBody(resp)
		if err != nil {
			if err != nil {
				statusCode := http.StatusInternalServerError
				statusText := http.StatusText(statusCode)
				w.WriteHeader(statusCode)
				w.Write([]byte(fmt.Sprintf("%s: %s", statusText, err)))
				return
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(respString))
}

func corsInterceptor(w http.ResponseWriter, r *http.Request, ah ApiHandler) {
	w.Header().Set("Access-Control-Allow-Origin", constants.CLIENT_BASE_URL)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(ah.Methods(), ", "))
	if r.Header.Get("Access-Control-Request-Headers") != "" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		status := http.StatusOK
		statusText := http.StatusText(http.StatusOK)
		w.WriteHeader(status)
		w.Write([]byte(statusText))
	}
}

func unmarshalRequestBody(req *http.Request, dest interface{}) error {
	bytes, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(bytes), dest)
}

func marshalResponseBody(resp interface{}) (string, error) {
	bytes, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
