package handlers

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/endpoints"
	"github.com/gorilla/mux"
)

func StartServer(port string, products map[int]database.Product) {
	router := mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	SetupHandlers(router, products)

	log.Infoln("Server running in port: ", port)
	http.ListenAndServe(port, router)
}

func SetupHandlers(router *mux.Router, products map[int]database.Product) {

	cartEndpoint := endpoints.NewCartEndpoint(products, os.Getenv("DISCOUNT_SERVICE_HOST"), os.Getenv("DISCOUNT_SERVICE_PORT"), os.Getenv("BLACK_FRIDAY"))

	router.HandleFunc("/cart", cartEndpoint.Post).Methods(http.MethodPost)
	log.Infoln("Endpoint /cart active.")

}
