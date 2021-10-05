// Copyright Christian Przybulinski
// All Rights Reserved

//Package endpoints used to setup the endpoints logic
package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
	log "github.com/sirupsen/logrus"
)

//CartRequest struct to Represent each Product Request from the Cart API
type CartRequest struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

//CartRequests struct the reprents a list of Product Request from the Cart API
type CartRequests struct {
	CartRequest []CartRequest `json:"products"`
}

//handleRequest is an Internal CartEndpoint method to handle the JSON request
//Transform the JSON into the struct CartRequests
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

	//validate if for some reason the cart is empty, so we can return an error
	if len(cartRequests.CartRequest) == 0 {
		log.Debugln("Cart is empty!!")
		err = errors.ErrEmptyCart
	} else {
		for _, c := range cartRequests.CartRequest {
			if c.ID < 1 || c.Quantity < 1 { //validate the ID and Quantity informed in the JSON request, so that it can't be lower than 1
				log.Debugln("ID or Quantity < 1")
				err = errors.ErrBadRequest
			}

		}
	}

	if err != nil {
		return CartRequests{}, err
	}
	return cartRequests, nil

}
