package ierrors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrType string

const (
	Unauthorized ErrType = "Unauthorized"
	Forbidden            = "Forbidden"
	BadRequest           = "BadRequest"
)

func NewError(errType ErrType, err error) error {
	return &Err{
		errType: errType,
		err:     err,
	}
}

type Err struct {
	errType ErrType
	err     error
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s: %s", e.errType, e.err)
}

func ToHttpStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var internalErr *Err
	if !errors.As(err, &internalErr) {
		return http.StatusInternalServerError
	}

	switch internalErr.errType {
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
