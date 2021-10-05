package endpoints

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

func TestCartEndpoint_sendResponse(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		response CartResponse
		err      error
	}

	cenario1 := args{mockResponseWriter(), CartResponse{}, nil}
	cenario2 := args{mockResponseWriter(), CartResponse{}, errors.ErrEmptyCart}

	tests := []struct {
		name string
		cart CartEndpoint
		args args
		want bool
	}{
		{"cenario 1", mockCartEndpoint(false), cenario1, true},
		{"cenario 2", mockCartEndpoint(false), cenario2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cart.sendResponse(tt.args.w, tt.args.response, tt.args.err); got != tt.want {
				t.Errorf("CartEndpoint.sendResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockResponseWriter() http.ResponseWriter {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ping")
	}

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	return w
}

func mockRequest(requestBody string) *http.Request {

	httpposturl := "localhost:8080/api/v1/cart"

	request, _ := http.NewRequest("POST", httpposturl, bytes.NewBuffer([]byte(requestBody)))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return request
}

func mockCartEndpoint(blackfriday bool) CartEndpoint {
	var cart CartEndpoint
	cart.Database = make(map[int]database.Product)
	cart.Database[1] = database.Product{Id: 1, Title: "Ergonomic Wooden Pants", Description: "Deleniti beatae porro.", Amount: 15157, Is_gift: false}
	cart.Database[2] = database.Product{Id: 2, Title: "Ergonomic Cotton Keyboard", Description: "Iste est ratione excepturi repellendus adipisci qui.", Amount: 93811, Is_gift: true}
	cart.Database[3] = database.Product{Id: 3, Title: "test", Description: "a little test.", Amount: 666, Is_gift: false}

	if blackfriday {
		cart.BlackFriday = time.Now().Format("2006-01-02")
	}

	return cart
}

func clearString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\n", ""), "\r", ""), "\t", "")
}

func getMockJSON(file string) string {
	var jsonResponse []byte

	jsonFile, err := os.Open(file)

	if err == nil {
		jsonResponse, err = ioutil.ReadAll(jsonFile)
	}

	if err == nil {
		return string(jsonResponse)
	}

	defer jsonFile.Close()
	return ""
}
