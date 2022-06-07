package error

import (
	"fmt"
	"net/http"
)

type RespError struct {
	Code    int
	Message string
	CodeMsg string
}

func (r *RespError) Error() string {
	return fmt.Sprintf("%d: %s", r.Code, r.Message)
}

func InternalServerError(msg string) error {
	return &RespError{
		Code:    http.StatusInternalServerError,
		Message: msg,
	}
}

func NotFound(msg string) error {
	return &RespError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func ValidationFailed(msg string) error {
	return &RespError{
		Code:    http.StatusConflict,
		Message: msg,
	}
}
