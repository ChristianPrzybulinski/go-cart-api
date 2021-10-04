package endpoints

import (
	"encoding/json"
	"math"
	"os"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/discount"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
	log "github.com/sirupsen/logrus"
)

type CartResponse struct {
	TotalAmount             int
	TotalAmountWithDiscount int
	TotalDiscount           int
	Products                []ResponseProduct
}

type ResponseProduct struct {
	ID          int
	Quantity    int
	UnitAmount  int
	TotalAmount int
	Discount    int
	IsGift      bool
}

func handleResponse(requests CartRequests, database map[int]database.Product) (CartResponse, error) {
	var response CartResponse
	var productMap map[int]ResponseProduct = make(map[int]ResponseProduct)

	response.TotalAmount = 0
	response.TotalAmountWithDiscount = 0
	response.TotalDiscount = 0

	for _, r := range requests.CartRequest {
		log.Debugln("Processing request ID: ", r.Id)

		rProduct, ok := handleProductRequest(r, database)
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
		if isBlackFriday(os.Getenv("BLACK_FRIDAY")) {
			rGift := getGift(database)
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

func handleProductRequest(r CartRequest, database map[int]database.Product) (ResponseProduct, bool) {

	if val, ok := database[r.Id]; ok {

		discountHost := os.Getenv("DISCOUNT_SERVICE_HOST")
		discountPort := os.Getenv("DISCOUNT_SERVICE_PORT")

		if len(discountHost) == 0 {
			discountHost = ":"
		}
		if len(discountPort) == 0 {
			discountPort = "50051"
		}

		dPercentage := discount.DescountPercentage(discountHost+discountPort, int32(r.Id))
		discountTotal := math.Round(float64(float32(val.Amount*r.Quantity) * dPercentage))
		return ResponseProduct{r.Id, r.Quantity, val.Amount, val.Amount * r.Quantity, int(discountTotal), false}, true
	}

	log.Errorln("Product ID not Found in database! Product not added..")
	return ResponseProduct{}, false
}

func (c CartResponse) JSON() (string, error) {
	res, err := json.Marshal(c)

	if err == nil {
		return string(res), nil
	} else {
		return "", err
	}

}
