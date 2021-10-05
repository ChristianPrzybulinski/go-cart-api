// Copyright Christian Przybulinski
// All Rights Reserved

//Package endpoints
package endpoints

import (
	"bytes"
	"encoding/json"
	"math"

	"github.com/ChristianPrzybulinski/go-cart-api/src/discount"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
	log "github.com/sirupsen/logrus"
)

//CartResponse struct to represent a Response from the Cart API
type CartResponse struct {
	TotalAmount             int               `json:"total_amount"`
	TotalAmountWithDiscount int               `json:"total_amount_with_discount"`
	TotalDiscount           int               `json:"total_discount"`
	Products                []ResponseProduct `json:"products"`
}

//ResponseProduct struct to represent the Product inside the response from the Cart API
type ResponseProduct struct {
	ID          int  `json:"id"`
	Quantity    int  `json:"quantity"`
	UnitAmount  int  `json:"unit_amount"`
	TotalAmount int  `json:"total_amount"`
	Discount    int  `json:"discount"`
	IsGift      bool `json:"is_gift"`
}

//handleResponse is an Internal CartEndpoint method that receives the Requests and will return the Response or error
func (cart CartEndpoint) handleResponse(requests CartRequests) (CartResponse, error) {
	var response CartResponse

	//sets everything to zero or empty
	var productMap map[int]ResponseProduct = make(map[int]ResponseProduct)
	response.TotalAmount = 0
	response.TotalAmountWithDiscount = 0
	response.TotalDiscount = 0

	//for each request it will sum the totals and add to the productMap
	for _, r := range requests.CartRequest {
		log.Debugln("Processing request ID: ", r.ID)

		rProduct, ok := cart.handleProductRequest(r)
		if ok {
			response.TotalAmount = response.TotalAmount + rProduct.TotalAmount
			response.TotalAmountWithDiscount = response.TotalAmountWithDiscount + (rProduct.TotalAmount - rProduct.Discount)
			response.TotalDiscount = response.TotalDiscount + rProduct.Discount

			if v, found := productMap[rProduct.ID]; found { //using a map structure in case the request comes with repeated ids in different positions
				v.TotalAmount = v.TotalAmount + rProduct.TotalAmount
				v.Quantity = v.Quantity + rProduct.Quantity
				v.Discount = v.Discount + rProduct.Discount
				productMap[rProduct.ID] = v
			} else {
				productMap[rProduct.ID] = rProduct
			}

			log.Debugln("Current Total Amount: ", response.TotalAmount)
			log.Debugln("Current Total Amount with Dicount: ", response.TotalAmount)
			log.Debugln("Current Total Discount: ", response.TotalDiscount)
		}
	}

	//after everything is done, turns the map back to slice, so we can add it in the response
	productSlice := mapToSlice(productMap)

	//before returning, it will check if we have any products processed in case we do, we need to check if its a black friday to add a gift
	if len(productMap) > 0 {
		if isBlackFriday(cart.BlackFriday) {
			rGift := getGift(cart.Database)
			if (rGift != ResponseProduct{}) {
				log.Debugln("Adding gift to the cart, BLACK FRIDAY baby")
				productSlice = append(productSlice, rGift)
			}
		}
		response.Products = productSlice
		return response, nil
	}

	return CartResponse{}, errors.ErrEmptyCart
}

//handleProductRequest is a submethod that handle one Product Request, returning it with the discount and maths already applied
func (cart CartEndpoint) handleProductRequest(r CartRequest) (ResponseProduct, bool) {

	if val, ok := cart.Database[r.ID]; ok {

		dPercentage := discount.DescountPercentage(cart.DiscountServiceAddress, int32(r.ID), cart.DiscountServiceTimeout)
		discountTotal := math.Round(float64(float32(val.Amount*r.Quantity) * dPercentage))
		return ResponseProduct{r.ID, r.Quantity, val.Amount, val.Amount * r.Quantity, int(discountTotal), false}, true
	}

	log.Errorln("Product ID not Found in database! Product not added..")
	return ResponseProduct{}, false
}

//JSON transforms the Response into a JSON
func (c CartResponse) JSON() string {
	res, _ := json.Marshal(c)
	var out bytes.Buffer

	json.Indent(&out, res, "", "  ")
	return out.String()

}
