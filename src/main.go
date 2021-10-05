package main

import (
	"os"
	"strings"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	"github.com/ChristianPrzybulinski/go-cart-api/src/handlers"
	"github.com/ChristianPrzybulinski/go-cart-api/src/utils"
	"github.com/mcuadros/go-defaults"
	log "github.com/sirupsen/logrus"
)

type Args struct {
	ApiAddress   string `default:":8080"`
	DatabasePath string `default:"./database"`
}

func Init() Args {
	var args Args
	defaults.SetDefaults(&args)

	if len(os.Args) <= 2 {
		utils.InitLog(strings.ToLower(os.Getenv("LOG_LEVEL")))
	} else {
		utils.InitLog(strings.ToLower(os.Args[2]))
	}

	if len(os.Args) <= 1 {
		host := os.Getenv("API_HOST")
		port := os.Getenv("API_PORT")

		if len(host) > 0 {
			if len(port) > 0 {
				args.ApiAddress = host + ":" + port
			} else {
				args.ApiAddress = host + args.ApiAddress
			}
		} else {
			if len(port) > 0 {
				args.ApiAddress = ":" + port
			}
		}
	} else {
		args.ApiAddress = os.Args[1]
	}

	if len(os.Args) <= 2 {
		db := os.Getenv("DATABASE_PATH")

		if len(db) > 0 {
			args.DatabasePath = db
		}
	} else {
		args.DatabasePath = os.Args[3]
	}

	return args
}

func main() {

	args := Init()

	log.Infoln("Logs started with " + log.GetLevel().String() + " level.")
	log.Infoln("Database path start with" + args.DatabasePath)

	log.Infoln("Loading Products Database...")
	products, err := database.GetAllProducts(args.DatabasePath + "/products.json")
	if err == nil {
		log.Infoln("Products Database Loaded.")

		handlers.StartServer(args.ApiAddress, products)
	} else {
		log.Errorln(err.Error())
		os.Exit(1)
	}
}
