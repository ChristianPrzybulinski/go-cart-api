// Copyright Christian Przybulinski
// All Rights Reserved

//Package handlers contains the handlers that are gonna be used in the API Server
package handlers

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/endpoints"
	"github.com/gorilla/mux"
)

//Args is used to define the API Address (Server:Port) and the DatabasePath since we are using internal json files as our database
type Args struct {
	APIAddress   string `default:":8080"`
	DatabasePath string `default:"./database"`
}

//StartServer method start the server, calling the handlers configured and listening in the address received inside the struct
func StartServer(args Args) {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter() //sets up a router with /api/v1/ as default path for all APIs
	SetupHandlers(router, args)

	log.Infoln("Server running in port: ", args.APIAddress)
	http.ListenAndServe(args.APIAddress, router)
}

//SetupHandlers defines the handlers that the API has, for now the only one is /cart
func SetupHandlers(router *mux.Router, args Args) {

	cartEndpoint := endpoints.NewCartEndpoint(args.DatabasePath,
		os.Getenv("DISCOUNT_SERVICE_HOST"),
		os.Getenv("DISCOUNT_SERVICE_PORT"),
		os.Getenv("DISCOUNT_SERVICE_TIMEOUT"),
		os.Getenv("BLACK_FRIDAY"))

	router.HandleFunc("/cart", cartEndpoint.Post).Methods(http.MethodPost)
	log.Infoln("Endpoint /cart active.")

}
