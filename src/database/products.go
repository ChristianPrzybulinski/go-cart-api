// Copyright Christian Przybulinski
// All Rights Reserved

//database package used to configure the JSONs files as database
package database

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

//Database Product, used to read the JSON file and convert to struct
type Product struct {
	ID          int    `json:"ID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	IsGift      bool   `json:"is_gift"`
}

//Read the JSON Database file and transform it to an internal Map, ID = key
func GetAllProducts(file string) (map[int]Product, error) {

	log.Debugln("Loading Json File: " + file)
	jsonFile, err := loadProducts(file)

	if err != nil {
		return nil, err
	}

	return jsonToMap(jsonFile)

}

//as the name says, transform the json received to a map of products
func jsonToMap(jsonFile []byte) (map[int]Product, error) {
	var products []Product

	log.Debugln("Json data: " + string(jsonFile))

	if len(jsonFile) > 0 {
		err := json.Unmarshal(jsonFile, &products)

		if err == nil {
			productMap := map[int]Product{}

			for _, product := range products {
				productMap[product.ID] = product
			}

			log.Debugln("Database size: ", len(productMap))
			return productMap, nil
		}
		return nil, err
	}

	return map[int]Product{}, nil

}

//read the json file and return its content
func loadProducts(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)

	if err == nil {
		return ioutil.ReadAll(jsonFile)
	}

	defer jsonFile.Close()

	return nil, err
}
