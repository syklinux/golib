package collection

import (
	"strings"
)

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	r1 := len(rs)
	start = start % r1
	if start < 0 {
		start = r1 + start
	}
	end := start + length
	if end > r1 {
		end = r1
	} else if end < start {
		start, end = end, start
		if start < 0 {
			start = 0
		}
	}
	return string(rs[start:end])
}

func ShortStr(str string, separator string, n int) string {
	strs := strings.Split(str, separator)
	// 缩短文件名，最多显示3级
	if n > len(strs) {
		n = len(strs)
	}
	result := ""
	for i := n; i > 0; i-- {
		result += strs[len(strs)-i] + separator
	}
	return strings.TrimSuffix(result, separator)
}
