package goutill

import (
	"crypto/md5"
	"fmt"
	"strings"
)

type stringUtil struct{}

var String = stringUtil{}

// Trim remove the space.
func (stringUtil) Trim(s string) string {
	var conv string
	convArr := strings.Split(s, " ")
	for _, v := range convArr {
		conv += v
	}

	return conv
}

// ToUpper changes the string to uppercase.
func (stringUtil) ToUpper(s string) string {
	return strings.ToUpper(String.Trim(s))
}

// Md5Hash changes the string to md5 hash string.
func (stringUtil) Md5Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
