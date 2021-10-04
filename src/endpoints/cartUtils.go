package endpoints

import (
	"time"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
	log "github.com/sirupsen/logrus"
)

func getGift(database map[int]database.Product) ResponseProduct {

	for _, r := range database {
		if r.Is_gift {
			return ResponseProduct{r.Id, 1, 0, 0, 0, true}
		}
	}
	log.Warnln("No Gift Product Found!")
	return ResponseProduct{}
}

func isBlackFriday(blackFriday string) bool {

	if len(blackFriday) == 0 {
		log.Debugln("Black Friday date not setted")
		return false
	} else {
		today := time.Now().Format("2006-01-02")
		return today == blackFriday
	}
}

func mapToSlice(mapProduct map[int]ResponseProduct) []ResponseProduct {
	var response []ResponseProduct

	for _, value := range mapProduct {
		response = append(response, value)
	}

	return response
}
