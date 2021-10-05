// Copyright Christian Przybulinski
// All Rights Reserved

//Package errors is used to set the HTTP Errors Response
package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	//ErrInternal is used as default error
	ErrInternal = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error (500)",
	}
	//ErrBadRequest is used when wrong Requests
	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400)",
	}
	//ErrNotFound is used as default not found page
	ErrNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Not Found (404)",
	}
	//ErrEmptyCart is used when the response is empty
	ErrEmptyCart = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400) - Empty Cart / no Product found!",
	}
)

//Error struct is used to encapsulate the error code and message to return in the json
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//GetError is used to get the specific http error or the default
func GetError(e error) *Error {
	res, ok := e.(*Error)
	if !ok {
		return ErrInternal
	}
	return res

}

//Error returns the error message with the code to log
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

//JSON returns the error in a JSON format
func (err *Error) JSON() string {
	var out bytes.Buffer
	if err == nil {
		return "{}"
	}
	res, _ := json.Marshal(err)

	json.Indent(&out, res, "", "  ")
	return out.String()
}

//StatusCode returns the status code of the error or OK in case its a not defined error
func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}
