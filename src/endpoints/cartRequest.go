package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
	log "github.com/sirupsen/logrus"
)

type CartRequest struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

type CartRequests struct {
	CartRequest []CartRequest `json:"products"`
}

func (cart CartEndpoint) handleRequest(r *http.Request) (CartRequests, error) {
	var cartRequests CartRequests
	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		log.Debugln("/cart API request: " + string(body))
		if len(string(body)) > 0 {
			err = json.Unmarshal(body, &cartRequests)
		} else {
			err = errors.ErrEmptyCart
		}
	}

	if len(cartRequests.CartRequest) == 0 {
		log.Debugln("Cart is empty!!")
		err = errors.ErrEmptyCart
	} else {
		for _, c := range cartRequests.CartRequest {
			if c.Id < 1 || c.Quantity < 1 {
				log.Debugln("ID or Quantity < 1")
				err = errors.ErrBadRequest
			}

		}
	}

	if err != nil {
		return CartRequests{}, err
	} else {
		return cartRequests, nil
	}

}
