package database

import (
	"reflect"
	"strconv"
	"testing"
)

type mocks struct {
	mockJson   string
	mockStruct []Product
}

func mockMaps(id []int, title []string, desc []string, amount []int, gift []bool) mocks {
	var products []Product

	jsonString := "["

	for i, _ := range id {
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
	}

	jsonString = jsonString + "]"

	return mocks{jsonString, products}
}

func mockMap(id int, title string, desc string, amount int, gift bool) mocks {

	product := []Product{{id, title, desc, amount, gift}}

	jsonString := "[{\"id\":" +
		strconv.Itoa(id) + ", \"title\": \"" +
		title + "\", \"description\": \"" +
		desc + "\", \"amount\": " +
		strconv.Itoa(amount) + ", \"is_gift\": " +
		strconv.FormatBool(gift) +
		"}]"

	return mocks{jsonString, product}
}

func mockMapWithError(id int, title string, desc string, amount int, gift bool) mocks {

	product := []Product{{id, title, desc, amount, gift}}

	jsonString := "[{\"id\":" +
		strconv.Itoa(id) + ", \"title\": \"" +
		title + "\", \"description\": " +
		strconv.Itoa(amount) + ", \"is_gift\": " +
		"}]"

	return mocks{jsonString, product}
}

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
		productMap[product.Id] = product
	}

	tests := []struct {
		name    string
		args    args
		want    map[int]Product
		wantErr bool
	}{
		{"working json", args{[]byte(mocksData[0].mockJson)},
			map[int]Product{mocksData[0].mockStruct[0].Id: mocksData[0].mockStruct[0]}, false},

		{"working json 2", args{[]byte(mocksData[1].mockJson)},
			map[int]Product{mocksData[1].mockStruct[0].Id: mocksData[1].mockStruct[0]}, false},

		{"working json 3", args{[]byte(mocksData[2].mockJson)}, productMap, false},

		{"working json vazio", args{[]byte("")}, map[int]Product{}, false},

		{"json incomplete 4", args{[]byte(mocksData[3].mockJson)},
			map[int]Product{mocksData[3].mockStruct[0].Id: mocksData[3].mockStruct[0]}, true},
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
		{"working case", args{"test.json"}, []byte("test"), false},
		{"not working case", args{"noexist.json"}, []byte("test"), true},
	}
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
