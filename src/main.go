// Copyright Christian Przybulinski
// All Rights Reserved

//Package main
package main

import (
	"os"
	"strings"

	"github.com/ChristianPrzybulinski/go-cart-api/src/handlers"
	"github.com/ChristianPrzybulinski/go-cart-api/src/utils"
	"github.com/mcuadros/go-defaults"
	log "github.com/sirupsen/logrus"
)

//Init used to init all the parameters and envvars used in the API
//Get args passed in the command line:
//1: API ADDRESS (host:port) ~ 2: Log Level (info,debug,warn,error) 3: Database Path (filesystem path)
func Init() handlers.Args {
	var args handlers.Args
	defaults.SetDefaults(&args)

	if len(os.Args) <= 2 {
		utils.InitLog(strings.ToLower(os.Getenv("LOG_LEVEL")))
	} else {
		utils.InitLog(strings.ToLower(os.Args[2]))
	}
	log.Infoln("Logs started with " + log.GetLevel().String() + " level.")

	if len(os.Args) <= 1 {
		host := os.Getenv("API_HOST")
		port := os.Getenv("API_PORT")

		if len(host) > 0 {
			if len(port) > 0 {
				args.APIAddress = host + ":" + port
			} else {
				args.APIAddress = host + args.APIAddress
			}
		} else {
			if len(port) > 0 {
				args.APIAddress = ":" + port
			}
		}
	} else {
		args.APIAddress = os.Args[1]
	}

	if len(os.Args) <= 3 {
		db := os.Getenv("DATABASE_PATH")

		if len(db) > 0 {
			args.DatabasePath = db
		}
	} else {
		args.DatabasePath = os.Args[3]
	}

	return args
}

//main method to start the server
func main() {

	args := Init()

	handlers.StartServer(args)
}
