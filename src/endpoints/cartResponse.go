package endpoints

import (
	"bytes"
	"encoding/json"
	"math"

	"github.com/ChristianPrzybulinski/go-cart-api/src/discount"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
	log "github.com/sirupsen/logrus"
)

type CartResponse struct {
	TotalAmount             int               `json:"total_amount"`
	TotalAmountWithDiscount int               `json:"total_amount_with_discount"`
	TotalDiscount           int               `json:"total_discount"`
	Products                []ResponseProduct `json:"products"`
}

type ResponseProduct struct {
	ID          int  `json:"id"`
	Quantity    int  `json:"quantity"`
	UnitAmount  int  `json:"unit_amount"`
	TotalAmount int  `json:"total_amount"`
	Discount    int  `json:"discount"`
	IsGift      bool `json:"is_gift"`
}

func (cart CartEndpoint) handleResponse(requests CartRequests) (CartResponse, error) {
	var response CartResponse
	var productMap map[int]ResponseProduct = make(map[int]ResponseProduct)

	response.TotalAmount = 0
	response.TotalAmountWithDiscount = 0
	response.TotalDiscount = 0

	for _, r := range requests.CartRequest {
		log.Debugln("Processing request ID: ", r.Id)

		rProduct, ok := cart.handleProductRequest(r)
		if ok {
			response.TotalAmount = response.TotalAmount + rProduct.TotalAmount
			response.TotalAmountWithDiscount = response.TotalAmountWithDiscount + (rProduct.TotalAmount - rProduct.Discount)
			response.TotalDiscount = response.TotalDiscount + rProduct.Discount

			if v, found := productMap[rProduct.ID]; found {
				v.TotalAmount = v.TotalAmount + rProduct.TotalAmount
				v.TotalAmount = v.Quantity + rProduct.Quantity
				v.TotalAmount = v.Discount + rProduct.Discount
			} else {
				productMap[rProduct.ID] = rProduct
			}

			log.Debugln("Current Total Amount: ", response.TotalAmount)
			log.Debugln("Current Total Amount with Dicount: ", response.TotalAmount)
			log.Debugln("Current Total Discount: ", response.TotalDiscount)
		}
	}

	productSlice := mapToSlice(productMap)

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

func (cart CartEndpoint) handleProductRequest(r CartRequest) (ResponseProduct, bool) {

	if val, ok := cart.Database[r.Id]; ok {

		dPercentage := discount.DescountPercentage(cart.DiscountServiceAddress, int32(r.Id))
		discountTotal := math.Round(float64(float32(val.Amount*r.Quantity) * dPercentage))
		return ResponseProduct{r.Id, r.Quantity, val.Amount, val.Amount * r.Quantity, int(discountTotal), false}, true
	}

	log.Errorln("Product ID not Found in database! Product not added..")
	return ResponseProduct{}, false
}

func (c CartResponse) JSON() string {
	res, _ := json.Marshal(c)
	var out bytes.Buffer

	json.Indent(&out, res, "", "  ")
	return out.String()

}
