package endpoints

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

type CartEndpoint struct {
	Database map[int]database.Product
}

func (cart CartEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	requests, err := handleRequest(r)

	if err == nil {
		r, err := handleResponse(requests, cart.Database)
		ok := cart.sendResponse(w, r, err)
		log.Infoln("The request was processed: ", ok)
	}
}

func (cart CartEndpoint) sendResponse(w http.ResponseWriter, response CartResponse, err error) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err == nil {
		responseJSON, err := response.JSON()
		if err == nil {
			fmt.Fprintf(w, responseJSON)
		}
	}

	w.WriteHeader(errors.GetError(err).StatusCode())

	log.Errorln(err.Error())
	fmt.Fprintf(w, errors.GetError(err).JSON())
	return false
}
