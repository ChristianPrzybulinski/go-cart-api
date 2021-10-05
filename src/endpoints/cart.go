package endpoints

import (
	"net/http"
	"strconv"

	"github.com/mcuadros/go-defaults"
	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

type CartEndpoint struct {
	Database               map[int]database.Product
	DiscountServiceAddress string `default:":50051"`
	DiscountServiceTimeout int    `default:"1"`
	BlackFriday            string `default:""`
}

func NewCartEndpoint(database map[int]database.Product, discountHost string, discountPort string, discountTimeout string, blackFriday string) Endpoint {
	var c CartEndpoint
	defaults.SetDefaults(&c)

	c.Database = database
	if len(discountHost) > 0 {
		if len(discountPort) > 0 {
			c.DiscountServiceAddress = discountHost + ":" + discountPort
		} else {
			temp := c.DiscountServiceAddress
			c.DiscountServiceAddress = discountHost + temp
		}
	} else {
		if len(discountPort) > 0 {
			c.DiscountServiceAddress = ":" + discountPort
		}
	}

	if len(discountTimeout) > 0 {
		timeout, err := strconv.Atoi(discountTimeout)
		if (err == nil) && (timeout > 0) {
			c.DiscountServiceTimeout = timeout
		}
	}

	if len(blackFriday) > 0 {
		c.BlackFriday = blackFriday
	}

	log.Infoln("Discount Service address: " + c.DiscountServiceAddress)
	log.Infoln("Discount Service Timeout: ", c.DiscountServiceTimeout)
	log.Infoln("Black Friday date: " + c.BlackFriday)

	return c
}

func (cart CartEndpoint) Post(w http.ResponseWriter, r *http.Request) {
	requests, err := cart.handleRequest(r)
	var ok bool
	if err == nil {
		r, err := cart.handleResponse(requests)
		ok = cart.sendResponse(w, r, err)

	} else {
		ok = cart.sendResponse(w, CartResponse{}, err)
	}
	log.Infoln("The request was processed: ", ok)
}

func (cart CartEndpoint) sendResponse(w http.ResponseWriter, response CartResponse, err error) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err == nil {
		responseJSON := response.JSON()
		log.Infoln("/cart API Response: ", responseJSON)
		w.Write([]byte(responseJSON))

		return true
	} else {
		log.Errorln(err.Error())
		w.WriteHeader(errors.GetError(err).StatusCode())
		responseJSON := errors.GetError(err).JSON()
		log.Infoln("/cart API Response: ", responseJSON)
		w.Write([]byte(responseJSON))
		return false
	}
}
