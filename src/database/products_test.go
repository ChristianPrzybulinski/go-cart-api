// Copyright Christian Przybulinski
// All Rights Reserved

package database

import (
	"reflect"
	"strconv"
	"testing"
)

//mocks struct used in the unit tests
type mocks struct {
	mockJSON   string
	mockStruct []Product
}

//mockMaps used to mock multiple single json item to map
func mockMaps(id []int, title []string, desc []string, amount []int, gift []bool) mocks {
	var products []Product

	jsonString := "["

	for i := range id {
		products = append(products, Product{id[i], title[i], desc[i], amount[i], gift[i]})

		if i > 0 {
			jsonString = jsonString + ", "
		}
		jsonString = jsonString + "{\"id\":" +
			strconv.Itoa(id[i]) + ", \"title\": \"" +
			title[i] + "\", \"description\": \"" +
			desc[i] + "\", \"amount\": " +
			strconv.Itoa(amount[i]) + ", \"is_gift\": " +
			strconv.FormatBool(gift[i]) +
			"}"
	} //since we need to test the read method too, we need to mock it as a string internally

	jsonString = jsonString + "]"

	return mocks{jsonString, products}
}

//mockMap used to mock one single json item to map
func mockMap(id int, title string, desc string, amount int, gift bool) mocks {

	product := []Product{{id, title, desc, amount, gift}}

	jsonString := "[{\"id\":" +
		strconv.Itoa(id) + ", \"title\": \"" +
		title + "\", \"description\": \"" +
		desc + "\", \"amount\": " +
		strconv.Itoa(amount) + ", \"is_gift\": " +
		strconv.FormatBool(gift) +
		"}]" //since we need to test the read method too, we need to mock it as a string internally

	return mocks{jsonString, product}
}

//mockMapWithError used to mock a database json with error
func mockMapWithError(id int, title string, desc string, amount int, gift bool) mocks {

	product := []Product{{id, title, desc, amount, gift}}

	jsonString := "[{\"id\":" +
		strconv.Itoa(id) + ", \"title\": \"" +
		title + "\", \"description\": " +
		strconv.Itoa(amount) + ", \"is_gift\": " +
		"}]" //since we need to test the read method too, we need to mock it as a string internally

	return mocks{jsonString, product}
}

//getMocks return the mocks to test
func getMocks() []mocks {

	mockData := []mocks{
		mockMap(1, "a", "b", 2, false),
		mockMap(132131, ",easd,asda,sd,as", "1203h812dn1", 432423, true),
		mockMaps([]int{1, 2}, []string{"a", "b"}, []string{"teste", "testeeee2"}, []int{3, 4}, []bool{false, true}),
		mockMapWithError(1, "a", "c", 2, false),
	}

	return mockData
}

func Test_jsonToMap(t *testing.T) {
	type args struct {
		jsonFile []byte
	}

	mocksData := getMocks()

	productMap := map[int]Product{}
	for _, product := range mocksData[2].mockStruct {
		productMap[product.ID] = product
	}

	tests := []struct {
		name    string
		args    args
		want    map[int]Product
		wantErr bool
	}{
		{"Simple JSON", args{[]byte(mocksData[0].mockJSON)},
			map[int]Product{mocksData[0].mockStruct[0].ID: mocksData[0].mockStruct[0]}, false},

		{"Simple JSON with different Data", args{[]byte(mocksData[1].mockJSON)},
			map[int]Product{mocksData[1].mockStruct[0].ID: mocksData[1].mockStruct[0]}, false},

		{"JSON with Multiple Data", args{[]byte(mocksData[2].mockJSON)}, productMap, false},

		{"Empty JSON", args{[]byte("")}, map[int]Product{}, false},

		{"Malformatted JSON", args{[]byte(mocksData[3].mockJSON)},
			map[int]Product{mocksData[3].mockStruct[0].ID: mocksData[3].mockStruct[0]}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonToMap(tt.args.jsonFile)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("jsonToMap() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("jsonToMap() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_loadProducts(t *testing.T) {
	type args struct {
		file string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Normal JSON", args{"unitTestData/3.json"}, []byte("test"), false},
		{"Wrong JSON path", args{"unitTestData/noexist.json"}, []byte("test"), true},
	} //using test files to check if the method actually reads their contents

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadProducts(tt.args.file)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("loadProducts() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("name = %v loadProducts() = %v, want %v", tt.name, got, tt.want)
				}
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	type args struct {
		file string
	}

	map1 := make(map[int]Product)
	map1[1] = Product{ID: 1, Title: "Ergonomic Wooden Pants", Description: "Deleniti beatae porro.", Amount: 15157, IsGift: false}

	map2 := make(map[int]Product)
	map2[1] = Product{ID: 1, Title: "Ergonomic Wooden Pants", Description: "Deleniti beatae porro.", Amount: 15157, IsGift: false}
	map2[2] = Product{ID: 2, Title: "Ergonomic Cotton Keyboard", Description: "Iste est ratione excepturi repellendus adipisci qui.", Amount: 93811, IsGift: true}

	tests := []struct {
		name    string
		args    args
		want    map[int]Product
		wantErr bool
	}{
		{"working case 1", args{"unitTestData/1.json"}, map1, false},
		{"working case 2", args{"unitTestData/2.json"}, map2, false},
		{"error case 1", args{"unitTestData/3.json"}, nil, true},
		{"error case 2", args{"unitTestData/nonexist.json"}, nil, true},
	} //since now we already have unit tests for each method separately, we can use some files to test the final method

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllProducts(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllProducts() = %v, want %v", got, tt.want)

			}
		})
	}
}
