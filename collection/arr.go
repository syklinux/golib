package collection

import (
	"strconv"
)

func StrArrUnique(array []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range array {
		if _, exist := keys[entry]; !exist {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func IntArrUnique(array []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range array {
		if _, exist := keys[entry]; !exist {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func StrArrToInt(strs []string) []int {
	result := make([]int, len(strs))
	for i, v := range strs {
		result[i], _ = strconv.Atoi(v)
	}
	return result
}

func StrArrDiff(arr1 []string, arrs ...[]string) (data []string) {
	if len(arrs) == 0 {
		return arr1
	}
	i := 0
loop:
	for {
		if i == len(arr1) {
			break
		}
		v := arr1[i]
		for _, arr := range arrs {
			for _, val := range arr {
				if v == val {
					i++
					continue loop
				}
			}
		}
		data = append(data, v)
		i++
	}
	return

}
