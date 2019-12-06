package goutill

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func TrimSpace(s string) string {
	var conv string
	convArr := strings.Split(s, " ")
	for _, v := range convArr {
		conv += v
	}

	return conv
}

func Normalization(s string) string {
	return strings.ToUpper(TrimSpace(s))
}

func Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}