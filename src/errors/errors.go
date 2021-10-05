// Copyright Christian Przybulinski
// All Rights Reserved

//Package used to set the HTTP Errors Response
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
	} //Used as default error
	ErrBadRequest = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400)",
	} //Used when wrong Requests
	ErrNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "Not Found (404)",
	} //Default not found page
	ErrEmptyCart = &Error{
		Code:    http.StatusBadRequest,
		Message: "Bad Request (400) - Empty Cart / no Product found!",
	} //Used when the response is empty
)

//struct to encapsulate the error code and message to return in the json
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//Method used to get the specific http error or the default
func GetError(e error) *Error {
	res, ok := e.(*Error)
	if !ok {
		return ErrInternal
	}
	return res

}

//return the error message with the code to log
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("error: code=%s message=%s", http.StatusText(err.Code), err.Message)
}

//return the error in a JSON format
func (err *Error) JSON() string {
	var out bytes.Buffer
	if err == nil {
		return "{}"
	}
	res, _ := json.Marshal(err)

	json.Indent(&out, res, "", "  ")
	return out.String()
}

//return the status code of the error or OK in case its a not defined error
func (err *Error) StatusCode() int {
	if err == nil {
		return http.StatusOK
	}
	return err.Code
}
