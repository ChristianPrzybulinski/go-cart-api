// Copyright Christian Przybulinski
// All Rights Reserved

//Package endpoints used to setup the endpoints logic
package endpoints

import (
	"net/http"
	"strconv"

	"github.com/mcuadros/go-defaults"
	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/errors"
)

//CartEndpoint struct that represents the CardEndpoints envvars and Database
type CartEndpoint struct {
	Database               map[int]database.Product
	DiscountServiceAddress string `default:":50051"`
	DiscountServiceTimeout int    `default:"1"`
	BlackFriday            string `default:""`
}

//NewCartEndpoint initialize the Endpoint, setting the envvars inside the struct or its default values
//For each parameters, its checked if its empty, in case true we use the default value
func NewCartEndpoint(databasePath string, discountHost string, discountPort string, discountTimeout string, blackFriday string) Endpoint {
	var c CartEndpoint
	defaults.SetDefaults(&c)

	log.Infoln("Database path is: " + databasePath)

	log.Infoln("Loading Products Database...")
	products, err := database.GetAllProducts(databasePath + "/products.json")

	if err != nil {
		log.Errorln(err.Error())
		c.Database = make(map[int]database.Product)
	} else {
		log.Infoln("Products Database Loaded.")
		c.Database = products
	}

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

//Post Method that handle the Post Request and Send the Response
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

//sendResponse writes the response back to the client, error or success (JSON)
func (cart CartEndpoint) sendResponse(w http.ResponseWriter, response CartResponse, err error) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err == nil {
		responseJSON := response.JSON()
		log.Infoln("/cart API Response: ", responseJSON)
		w.Write([]byte(responseJSON))
		return true
	}

	log.Errorln(err.Error())
	w.WriteHeader(errors.GetError(err).StatusCode())
	responseJSON := errors.GetError(err).JSON()
	log.Infoln("/cart API Response: ", responseJSON)
	w.Write([]byte(responseJSON))
	return false

}
