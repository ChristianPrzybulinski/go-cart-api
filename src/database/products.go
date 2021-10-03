package database

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type Product struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Is_gift     bool   `json:"is_gift"`
}

func GetAllProducts(file string) (map[int]Product, error) {

	log.Debugln("Loading Json File: " + file)
	jsonFile, err := loadProducts(file)

	if err != nil {
		return nil, err
	} else {
		return jsonToMap(jsonFile)
	}

}

func jsonToMap(jsonFile []byte) (map[int]Product, error) {
	var products []Product

	log.Debugln("Json data: " + string(jsonFile))

	if len(jsonFile) > 0 {
		err := json.Unmarshal(jsonFile, &products)

		if err == nil {
			productMap := map[int]Product{}

			for _, product := range products {
				productMap[product.Id] = product
			}

			log.Debugln("Database size: ", len(productMap))
			return productMap, nil
		}
		return nil, err
	} else {
		return map[int]Product{}, nil
	}

}

func loadProducts(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)

	if err == nil {
		return ioutil.ReadAll(jsonFile)
	}

	defer jsonFile.Close()

	return nil, err
}
