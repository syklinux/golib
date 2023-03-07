package utils

import (
	"fmt"
)

import "crypto/md5"

func Md5(str string) string {
	key := []byte(str)
	hash := md5.Sum(key)
	return fmt.Sprintf("%x", hash)
}
