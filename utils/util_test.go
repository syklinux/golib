package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetStringSliceMapValue(t *testing.T) {
	temp := map[string]interface{}{
		"city_cvr_sort": "[\"5e72f32f315a77f7298b8937\"]",
	}
	test := Convert2StringSlice(temp["city_cvr_sort"])
	fmt.Println(test)
}

func Convert2StringSlice(data interface{}) []string {
	if data == nil {
		return nil
	}
	dataStr, ok := data.(string)
	if !ok {
		return nil
	}
	var valMap []string
	err := json.Unmarshal([]byte(dataStr), &valMap)
	if err != nil {
		return nil
	}
	return valMap
}

func TestPageSizeV2(t *testing.T) {
	pList := []string{
		"1", "2", "3", "4", "5", "6",
	}
	pList = nil
	page := 2
	size := 2

	a, b := PageSizeV2(len(pList), page, size)
	i := PageSize(pList, page, size)
	fmt.Printf("list:%v,", i)
	fmt.Printf("a:%v,b:%v,", a, b)
	fmt.Printf("list:%v", pList[a:b])
}

func TestContains(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if len(list) > 5 {
		list = list[:5]
	}
	fmt.Printf("%v", list)
}

func TestCalEarthDistance(t *testing.T) {
	version := "9.10.123"
	newVersion := "9.11.0"
	flag := "<="
	if version > newVersion {
		flag = ">"
	}
	fmt.Printf("version:%v %v newVersion:%v", version, flag, newVersion)

}

func TestContainsInt(t *testing.T) {
	temp := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 10, 11}

	if len(temp) > 9 {
		temp = temp[:9]
	}
	fmt.Printf("%v", temp)
}
