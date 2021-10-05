package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error (500)",
	}
	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400)",
	}
	ErrNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Not Found (404)",
	}
	ErrEmptyCart = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400) - Empty Cart / no Product found!",
	}
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetError(e error) *Error {
	res, ok := e.(*Error)
	if !ok {
		return ErrInternal
	} else {
		return res
	}
}

func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

func (err *Error) JSON() string {
	var out bytes.Buffer
	if err == nil {
		return "{}"
	}
	res, _ := json.Marshal(err)

	json.Indent(&out, res, "", "  ")
	return out.String()
}

func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}
