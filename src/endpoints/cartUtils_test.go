package endpoints

import (
	"reflect"
	"testing"
	"time"

	"github.com/ChristianPrzybulinski/go-cart-api/src/database"
)

func Test_isBlackFriday(t *testing.T) {
	type args struct {
		blackFriday string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Black Friday Day", args{time.Now().Format("2006-01-02")}, true},
		{"Not Black Friday Day", args{"2001-23-11"}, false},
		{"Empty string", args{""}, false},
		{"Wrong date string", args{"22222-333-12"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBlackFriday(tt.args.blackFriday); got != tt.want {
				t.Errorf("isBlackFriday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGift(t *testing.T) {
	type args struct {
		database map[int]database.Product
	}

	map1 := make(map[int]database.Product)
	map1[0] = database.Product{Id: 0, Title: "0", Description: "0", Amount: 13123, Is_gift: true}

	response1 := []ResponseProduct{
		{0, 1, 0, 0, 0, true},
	}

	map2 := make(map[int]database.Product)
	map2[0] = database.Product{Id: 0, Title: "123", Description: "asda", Amount: 11, Is_gift: false}
	map2[1] = database.Product{Id: 1, Title: "0", Description: "0", Amount: 13123, Is_gift: true}

	response2 := []ResponseProduct{
		{1, 1, 0, 0, 0, true},
	}

	map3 := make(map[int]database.Product)
	map3[0] = database.Product{Id: 0, Title: "123", Description: "asda", Amount: 11, Is_gift: true}
	map3[1] = database.Product{Id: 1, Title: "0", Description: "0", Amount: 13123, Is_gift: true}

	response3 := []ResponseProduct{
		{0, 1, 0, 0, 0, true},
		{1, 1, 0, 0, 0, true},
	}

	map4 := make(map[int]database.Product)
	map4[0] = database.Product{Id: 0, Title: "123", Description: "asda", Amount: 11, Is_gift: false}
	map4[1] = database.Product{Id: 1, Title: "0", Description: "0", Amount: 13123, Is_gift: false}

	response4 := []ResponseProduct{
		{},
	}

	tests := []struct {
		name              string
		args              args
		want              []ResponseProduct
		moreThanOneResult bool
	}{
		{"Gift found", args{map1}, response1, false},
		{"Gift found with two value", args{map2}, response2, false},
		{"Two values as Gifts found", args{map3}, response3, true},
		{"no Gifts found", args{map4}, response4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.moreThanOneResult {
				var ok bool
				ok = false

				for i := range tt.want {
					if got := getGift(tt.args.database); reflect.DeepEqual(got, tt.want[i]) {
						ok = true
					}
				}
				if !ok {
					t.Errorf("getGift() = %v, want %v", "None of the elements were gifts", tt.want[0])
				}

			} else {
				if got := getGift(tt.args.database); !reflect.DeepEqual(got, tt.want[0]) {
					t.Errorf("getGift() = %v, want %v", got, tt.want[0])
				}
			}

		})
	}
}
