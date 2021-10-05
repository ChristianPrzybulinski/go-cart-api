// Copyright Christian Przybulinski
// All Rights Reserved

package errors

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestGetError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{"Other Errors", args{http.ErrBodyNotAllowed}, ErrInternal},
		{"ErrNotFound", args{ErrNotFound}, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetError(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{"ErrNotFound", ErrNotFound, "error: code=" + http.StatusText(http.StatusNotFound) + " message=Not Found (404)"},
		{"StatusInternalServerError", ErrInternal, "error: code=" + http.StatusText(http.StatusInternalServerError) + " message=Internal Server Error (500)"},
		{"StatusBadRequest", ErrBadRequest, "error: code=" + http.StatusText(http.StatusBadRequest) + " message=Bad Request (400)"},
		{"Nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_JSON(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{"ErrNotFound", ErrNotFound, "{\"code\":404,\"message\":\"Not Found (404)\"}"},
		{"StatusInternalServerError", ErrInternal, "{\"code\":500,\"message\":\"Internal Server Error (500)\"}"},
		{"StatusBadRequest", ErrBadRequest, "{\"code\":400,\"message\":\"Bad Request (400)\"}"},
		{"Nil", nil, "{}"},
		{"Empty Cart", ErrEmptyCart, "{\"code\":400,\"message\":\"Bad Request (400) - Empty Cart / no Product found!\"}"},
	} //could be used some files to test, but since its just a string to JSON method, not necessary
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.JSON(); clearString(got) != clearString(tt.want) {
				t.Errorf("Error.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_StatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want int
	}{
		{"ErrNotFound", ErrNotFound, http.StatusNotFound},
		{"StatusInternalServerError", ErrInternal, http.StatusInternalServerError},
		{"StatusBadRequest", ErrBadRequest, http.StatusBadRequest},
		{"Nil", nil, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.StatusCode(); got != tt.want {
				t.Errorf("Error.StatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

//clearString utils method to clear all spaces and newlines to compare
func clearString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\n", ""), "\r", ""), "\t", "")
}
