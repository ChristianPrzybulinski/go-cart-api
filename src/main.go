package main

import (
	"os"
	"strings"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/handlers"
	"github.com/ChristianPrzybulinski/go-cart-api/src/utils"
	log "github.com/sirupsen/logrus"
)

func main() {

	utils.InitLog(strings.ToLower(os.Getenv("LOG_LEVEL")))
	log.Infoln("Logs started with " + log.GetLevel().String() + " level.")

	log.Infoln("Loading Products Database...")

	databasePath := os.Getenv("DATABASE_PATH")

	if len(databasePath) == 0 {
		databasePath = "./database"
	}

	products, err := database.GetAllProducts(databasePath + "/products.json")
	if err == nil {
		log.Infoln("Products Database Loaded.")

		hostname := os.Getenv("API_HOST")
		hostport := os.Getenv("API_PORT")

		if len(hostname) == 0 {
			hostname = ":"
		}
		if len(hostport) == 0 {
			hostport = "8080"
		}

		handlers.StartServer(hostname+hostport, products)
	} else {
		log.Errorln(err.Error())
		os.Exit(1)
	}
}
