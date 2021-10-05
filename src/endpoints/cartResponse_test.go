package endpoints

import (
	"reflect"
	"testing"
)

func TestCartResponse_JSON(t *testing.T) {
	tests := []struct {
		name    string
		c       CartResponse
		want    string
		wantErr bool
	}{
		{"working case 1", CartResponse{TotalAmount: 100, TotalAmountWithDiscount: 50, TotalDiscount: 50,
			Products: []ResponseProduct{{ID: 1, Quantity: 1, UnitAmount: 100, TotalAmount: 100, Discount: 50, IsGift: false}}}, getMockJSON("unitTestData/responses/3.json"), false},
		{"working case 2", CartResponse{TotalAmount: 200, TotalAmountWithDiscount: 150, TotalDiscount: 50,
			Products: []ResponseProduct{{ID: 1, Quantity: 1, UnitAmount: 100, TotalAmount: 100, Discount: 50, IsGift: false},
				{ID: 2, Quantity: 1, UnitAmount: 100, TotalAmount: 100, Discount: 0, IsGift: true}}}, getMockJSON("unitTestData/responses/2.json"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.JSON()

			if clearString(got) != clearString(tt.want) {
				t.Errorf("CartResponse.JSON() = %v, want %v", clearString(got), clearString(tt.want))
			}
		})
	}
}

func TestCartEndpoint_handleProductRequest(t *testing.T) {
	type args struct {
		r CartRequest
	}

	request1 := CartRequest{Id: 1, Quantity: 1}
	response1 := ResponseProduct{ID: 1, Quantity: 1, UnitAmount: 15157, TotalAmount: 15157, IsGift: false}

	request2 := CartRequest{Id: 2, Quantity: 4}
	response2 := ResponseProduct{ID: 2, Quantity: 4, UnitAmount: 93811, TotalAmount: 375244, IsGift: false}

	request3 := CartRequest{Id: 3, Quantity: 2}
	response3 := ResponseProduct{ID: 3, Quantity: 2, UnitAmount: 666, TotalAmount: 1332, IsGift: false}

	request4 := CartRequest{Id: 4, Quantity: 1}

	tests := []struct {
		name  string
		cart  CartEndpoint
		args  args
		want  ResponseProduct
		want1 bool
	}{
		{"working case", mockCartEndpoint(false), args{request1}, response1, true},
		{"working case 2", mockCartEndpoint(false), args{request2}, response2, true},
		{"working case 3", mockCartEndpoint(true), args{request3}, response3, true},
		{"not found", mockCartEndpoint(true), args{request4}, ResponseProduct{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.cart.handleProductRequest(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartEndpoint.handleProductRequest() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CartEndpoint.handleProductRequest() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCartEndpoint_handleResponse(t *testing.T) {
	type args struct {
		requests CartRequests
	}

	cartrequests1 := CartRequests{CartRequest: []CartRequest{
		{Id: 1, Quantity: 1}}}
	cartresponse1 := CartResponse{TotalAmount: 15157, TotalAmountWithDiscount: 15157, TotalDiscount: 0,
		Products: []ResponseProduct{
			{ID: 1, Quantity: 1, UnitAmount: 15157, TotalAmount: 15157, IsGift: false}}}

	cartrequests2 := CartRequests{CartRequest: []CartRequest{
		{Id: 1, Quantity: 1},
		{Id: 2, Quantity: 4}}}
	cartresponse2 := CartResponse{TotalAmount: 390401, TotalAmountWithDiscount: 390401, TotalDiscount: 0,
		Products: []ResponseProduct{
			{ID: 1, Quantity: 1, UnitAmount: 15157, TotalAmount: 15157, IsGift: false},
			{ID: 2, Quantity: 4, UnitAmount: 93811, TotalAmount: 375244, IsGift: false}}}

	cartrequests3 := CartRequests{CartRequest: []CartRequest{
		{Id: 1, Quantity: 1},
		{Id: 2, Quantity: 4}}}
	cartresponse3 := CartResponse{TotalAmount: 390401, TotalAmountWithDiscount: 390401, TotalDiscount: 0,
		Products: []ResponseProduct{
			{ID: 1, Quantity: 1, UnitAmount: 15157, TotalAmount: 15157, IsGift: false},
			{ID: 2, Quantity: 4, UnitAmount: 93811, TotalAmount: 375244, IsGift: false},
			{ID: 2, Quantity: 1, UnitAmount: 0, TotalAmount: 0, IsGift: true}}}

	cartrequests4 := CartRequests{CartRequest: []CartRequest{
		{Id: 4, Quantity: 1}}}

	cartrequests5 := CartRequests{CartRequest: []CartRequest{
		{Id: 2, Quantity: 4},
		{Id: 1, Quantity: 1},
		{Id: 2, Quantity: 4},
	}}
	cartresponse5 := CartResponse{TotalAmount: 765645, TotalAmountWithDiscount: 765645, TotalDiscount: 0,
		Products: []ResponseProduct{
			{ID: 1, Quantity: 1, UnitAmount: 15157, TotalAmount: 15157, IsGift: false},
			{ID: 2, Quantity: 8, UnitAmount: 93811, TotalAmount: 750488, IsGift: false}}}

	tests := []struct {
		name    string
		cart    CartEndpoint
		args    args
		want    CartResponse
		wantErr bool
	}{
		{"working case 1", mockCartEndpoint(false), args{cartrequests1}, cartresponse1, false},
		{"working case 2", mockCartEndpoint(false), args{cartrequests2}, cartresponse2, false},
		{"black friday", mockCartEndpoint(true), args{cartrequests3}, cartresponse3, false},
		{"error", mockCartEndpoint(true), args{cartrequests4}, CartResponse{}, true},
		{"using map strategy", mockCartEndpoint(false), args{cartrequests5}, cartresponse5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cart.handleResponse(tt.args.requests)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartEndpoint.handleResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartEndpoint.handleResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
