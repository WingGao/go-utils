package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5Sum(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
