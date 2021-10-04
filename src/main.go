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
	products, err := database.GetAllProducts(os.Getenv("DATABASE_PATH") + "/products.json")
	if err == nil {
		log.Infoln("Products Database Loaded.")
	} else {
		log.Errorln(err.Error())
	}
	handlers.StartServer(":"+os.Getenv("API_PORT"), products)

}
