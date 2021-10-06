// Copyright Christian Przybulinski
// All Rights Reserved

package endpoints

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

func TestNewCartEndpoint(t *testing.T) {
	type args struct {
		database        string
		discountHost    string
		discountPort    string
		discountTimeout string
		blackFriday     string
	}

	dbTest := make(map[int]database.Product)
	dbTest[1] = database.Product{ID: 1, Title: "Ergonomic Wooden Pants", Description: "Deleniti beatae porro.", Amount: 15157, IsGift: false}

	tests := []struct {
		name string
		args args
		want Endpoint
	}{
		{"All parameters informed", args{"", "test", "123", "2", "2021-10-23"},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: "test:123", DiscountServiceTimeout: 2, BlackFriday: "2021-10-23"}},

		{"Without black friday", args{"", "", "444", "50", ""},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: ":444", DiscountServiceTimeout: 50, BlackFriday: ""}},

		{"without port, blackfriday and timeout", args{"", "testing", "", "", ""},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: "testing:50051", DiscountServiceTimeout: 1, BlackFriday: ""}},

		{"only database informed", args{"unitTestData/databases", "", "", "0", ""},
			CartEndpoint{Database: dbTest, DiscountServiceAddress: ":50051", DiscountServiceTimeout: 1, BlackFriday: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCartEndpoint(tt.args.database, tt.args.discountHost, tt.args.discountPort, tt.args.discountTimeout, tt.args.blackFriday); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCartEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{"Success scenario", mockCartEndpoint(false), cenario1, true},
		{"Failed scenario", mockCartEndpoint(false), cenario2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cart.sendResponse(tt.args.w, tt.args.response, tt.args.err); got != tt.want {
				t.Errorf("CartEndpoint.sendResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

//mockResponseWriter mock a responseWriter to use
func mockResponseWriter() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	return w
}

//mockRequest mock a HttpRequest to use
func mockRequest(requestBody string) *http.Request {

	req := httptest.NewRequest("POST", "http://example.com/foo", strings.NewReader(requestBody))

	return req
}

//mockCartEndpoint mock a CartEndpoint with some default database values and passing if we want black friday or not
func mockCartEndpoint(blackfriday bool) CartEndpoint {
	var cart CartEndpoint
	cart.Database = make(map[int]database.Product)
	cart.Database[1] = database.Product{ID: 1, Title: "Ergonomic Wooden Pants", Description: "Deleniti beatae porro.", Amount: 15157, IsGift: false}
	cart.Database[2] = database.Product{ID: 2, Title: "Ergonomic Cotton Keyboard", Description: "Iste est ratione excepturi repellendus adipisci qui.", Amount: 93811, IsGift: true}
	cart.Database[3] = database.Product{ID: 3, Title: "test", Description: "a little test.", Amount: 666, IsGift: false}

	if blackfriday {
		cart.BlackFriday = time.Now().Format("2006-01-02")
	}

	return cart
}

//clearString clear the whitespaces and newlines to compare
func clearString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\n", ""), "\r", ""), "\t", "")
}

//getMockJSON reads the file (path) informed and return it as a string (used to read the file that contains the json content)
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

func TestCartEndpoint_Post(t *testing.T) {
	type args struct {
		req  string
		want string
	}
	tests := []struct {
		name string
		cart CartEndpoint
		args args
	}{
		{"Using a mocked database", mockCartEndpoint(false), args{getMockJSON("unitTestData/requests/1.json"), getMockJSON("unitTestData/responses/1.json")}},
		{"using a mocked database with different data", mockCartEndpoint(false), args{"unitTestData/requests/4.json", getMockJSON("unitTestData/responses/4.json")}},
		{"real database without black friday", mockRealCartEndpoint("unitTestData/databases/1.json", false), args{getMockJSON("unitTestData/requests/5.json"), getMockJSON("unitTestData/responses/5.json")}},
		{"empty json request", mockRealCartEndpoint("unitTestData/databases/1.json", false), args{getMockJSON("unitTestData/requests/4.json"), getMockJSON("unitTestData/responses/4.json")}},
		{"real database with black friday", mockRealCartEndpoint("unitTestData/databases/1.json", true), args{getMockJSON("unitTestData/requests/5.json"), getMockJSON("unitTestData/responses/6.json")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w, r := mockServer(tt.args.req)
			tt.cart.Post(w, r)

			resp := w.Result()
			got, _ := ioutil.ReadAll(resp.Body)

			if !reflect.DeepEqual(clearString(string(got)), clearString(tt.args.want)) {
				t.Errorf("TestCartEndpoint_Post() = %v, want %v", clearString(string(got)), clearString(tt.args.want))
			}
		})
	}
}

//mockServer mocks a response and request to use
func mockServer(request string) (*httptest.ResponseRecorder, *http.Request) {
	req := mockRequest(request)
	w := mockResponseWriter()
	return w, req
}

//mockRealCartEndpoint mocks a CartEndpoint with a real database file json, and using blackfriday true or false as a parameter
func mockRealCartEndpoint(file string, blackfriday bool) CartEndpoint {
	var cart CartEndpoint
	m, _ := database.GetAllProducts(file)

	cart.Database = m

	if blackfriday {
		cart.BlackFriday = time.Now().Format("2006-01-02")
	}

	return cart
}
