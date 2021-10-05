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
		database        map[int]database.Product
		discountHost    string
		discountPort    string
		discountTimeout string
		blackFriday     string
	}
	tests := []struct {
		name string
		args args
		want Endpoint
	}{
		{"test 1", args{make(map[int]database.Product), "test", "123", "2", "2021-10-23"},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: "test:123", DiscountServiceTimeout: 2, BlackFriday: "2021-10-23"}},

		{"test 2", args{make(map[int]database.Product), "", "444", "50", ""},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: ":444", DiscountServiceTimeout: 50, BlackFriday: ""}},

		{"test 3", args{make(map[int]database.Product), "testing", "", "", ""},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: "testing:50051", DiscountServiceTimeout: 1, BlackFriday: ""}},

		{"test 4", args{make(map[int]database.Product), "", "", "0", ""},
			CartEndpoint{Database: make(map[int]database.Product), DiscountServiceAddress: ":50051", DiscountServiceTimeout: 1, BlackFriday: ""}},
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

func mockResponseWriter() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()

	return w
}

func mockRequest(requestBody string) *http.Request {

	req := httptest.NewRequest("POST", "http://example.com/foo", strings.NewReader(requestBody))

	return req
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
		{"test 1", mockCartEndpoint(false), args{getMockJSON("unitTestData/requests/1.json"), getMockJSON("unitTestData/responses/1.json")}},
		{"test 2", mockCartEndpoint(false), args{"unitTestData/requests/4.json", getMockJSON("unitTestData/responses/4.json")}},
		{"test 3", mockRealCartEndpoint("unitTestData/databases/1.json", false), args{getMockJSON("unitTestData/requests/5.json"), getMockJSON("unitTestData/responses/5.json")}},
		{"test 4", mockRealCartEndpoint("unitTestData/databases/1.json", false), args{getMockJSON("unitTestData/requests/4.json"), getMockJSON("unitTestData/responses/4.json")}},
		{"test 5", mockRealCartEndpoint("unitTestData/databases/1.json", true), args{getMockJSON("unitTestData/requests/5.json"), getMockJSON("unitTestData/responses/6.json")}},
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

func mockServer(request string) (*httptest.ResponseRecorder, *http.Request) {
	req := mockRequest(request)
	w := mockResponseWriter()
	return w, req
}

func mockRealCartEndpoint(file string, blackfriday bool) CartEndpoint {
	var cart CartEndpoint
	m, _ := database.GetAllProducts(file)

	cart.Database = m

	if blackfriday {
		cart.BlackFriday = time.Now().Format("2006-01-02")
	}

	return cart
}
