package main

import (
	"os"

	"github.com/ChristianPrzybulinski/go-cart-api/src/handlers"
)

func main() {

	handlers.StartServer(":" + os.Getenv("API_PORT"))
}
