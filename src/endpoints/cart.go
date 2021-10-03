package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

type CartEndpoint struct {
	Database map[int]database.Product
}

type CartRequests struct {
	CartRequest []struct {
		Id       int `json:"id"`
		Quantity int `json:"quantity"`
	} `json:"products"`
}

func (cart CartEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	requests, err := cart.handleRequest(r)

	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		cart.sendResponse(w, requests)
	} else {
		log.Errorln(err.Error())
		fmt.Fprintf(w, errors.GetError(err).JSON()) //to-do validar se eh o melhor metodo para retornar
	}
}

func (cart CartEndpoint) sendResponse(w http.ResponseWriter, requests CartRequests) error {
	fmt.Fprintf(w, cart.Database[requests.CartRequest[0].Id].Description)
	return nil
}

func (cart CartEndpoint) handleRequest(r *http.Request) (CartRequests, error) {
	var cartRequests CartRequests
	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		log.Debugln("/cart API request: " + string(body))
		err = json.Unmarshal(body, &cartRequests)
	}

	return cartRequests, err
}
