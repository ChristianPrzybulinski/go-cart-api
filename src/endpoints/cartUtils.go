// Copyright Christian Przybulinski
// All Rights Reserved

//Endpoints package
package endpoints

import (
	"sort"
	"time"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	log "github.com/sirupsen/logrus"
)

//Used to get a Gift from the database, it will return as default the first one that it finds (no order implemented)
func getGift(database map[int]database.Product) ResponseProduct {

	for _, r := range database {
		if r.IsGift {
			return ResponseProduct{r.ID, 1, 0, 0, 0, true}
		}
	}
	log.Warnln("No Gift Product Found!")
	return ResponseProduct{}
}

//Check if the date received is today, the date format is YYYY-MM-DD
func isBlackFriday(blackFriday string) bool {

	if len(blackFriday) == 0 {
		log.Debugln("Black Friday date not setted")
		return false
	}
	today := time.Now().Format("2006-01-02")
	return today == blackFriday

}

//Transform a Map back to a Slice and order ascending by ID
func mapToSlice(mapProduct map[int]ResponseProduct) []ResponseProduct {
	var response []ResponseProduct

	for _, value := range mapProduct {
		response = append(response, value)
	}

	sort.SliceStable(response, func(i, j int) bool {
		return response[i].ID < response[j].ID
	})

	return response
}
