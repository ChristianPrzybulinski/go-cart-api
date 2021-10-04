package endpoints

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func Test_handleRequest(t *testing.T) {
	type args struct {
		r *http.Request
	}

	request1 := CartRequests{[]CartRequest{{1, 1}}}
	request2 := CartRequests{[]CartRequest{{1, 1}, {4, 1231}}}

	tests := []struct {
		name    string
		args    args
		want    CartRequests
		wantErr bool
	}{
		{"Valid request", args{mockRequest("{\"products\": [{\"id\": 1,\"quantity\": 1 }]}")}, request1, false},
		{"invalid request", args{mockRequest("{\"saaaa\": \"id\": 1,\"quantity\": 1 }]}")}, CartRequests{}, true},
		{"invalid request 2", args{mockRequest("{\"asdsadasdas\": [{\"id\": 1,\"quantity\": 1 }]}")}, CartRequests{}, true},
		{"invalid request 3", args{mockRequest("{\"products\": [{\"aaa\": 1,\"quantity\": 1 }]}")}, CartRequests{}, true},
		{"invalid request 3", args{mockRequest("{\"products\": [{\"id\": 1,\"quantity\": asda }]}")}, CartRequests{}, true},
		{"Double requests", args{mockRequest("{\"products\": [{\"id\": 1,\"quantity\": 1 }, {\"id\": 4,\"quantity\": 1231 }]}")}, request2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleRequest(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleRequest() = %v, want %v", got, tt.want)
			}

		})
	}
}

func mockRequest(requestBody string) *http.Request {

	httpposturl := "localhost:8080/api/v1/cart"

	request, _ := http.NewRequest("POST", httpposturl, bytes.NewBuffer([]byte(requestBody)))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return request
}
