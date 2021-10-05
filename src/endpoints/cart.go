package endpoints

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

type CartEndpoint struct {
	Database               map[int]database.Product
	DiscountServiceAddress string `default:":50051"`
	BlackFriday            string `default:""`
}

func NewCartEndpoint(database map[int]database.Product, discountHost string, discountPort string, blackFriday string) Endpoint {
	var c CartEndpoint

	c.Database = database
	if len(discountHost) > 0 {
		if len(discountPort) > 0 {
			c.DiscountServiceAddress = discountHost + ":" + discountPort
		} else {
			c.DiscountServiceAddress = discountHost + c.DiscountServiceAddress
		}
	} else {
		if len(discountPort) > 0 {
			c.DiscountServiceAddress = ":" + discountPort
		}
	}

	if len(blackFriday) > 0 {
		c.BlackFriday = blackFriday
	}

	log.Infoln("Discount Service address: " + c.DiscountServiceAddress)
	log.Infoln("Black Friday date: " + c.BlackFriday)

	return c
}

func (cart CartEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	requests, err := cart.handleRequest(r)

	if err == nil {
		r, err := cart.handleResponse(requests)
		ok := cart.sendResponse(w, r, err)
		log.Infoln("The request was processed: ", ok)
	}
}

func (cart CartEndpoint) sendResponse(w http.ResponseWriter, response CartResponse, err error) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err == nil {
		w.Write([]byte(response.JSON()))
		return true
	} else {
		log.Errorln(err.Error())
		w.WriteHeader(errors.GetError(err).StatusCode())
		w.Write([]byte(errors.GetError(err).JSON()))
		return false
	}
}
