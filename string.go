package utils

import (
	"math/rand"
	"time"
	"strconv"
	"strings"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("0123456789")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandIntString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}
	return string(b)
}

func SubString(s string, start, length int) (sub string) {
	length = MinInt(length, len(s))
	sr := []rune(s) // for unicode
	return string(sr[start:length])
}

func StrIsEmpty(s string) bool {
	return s == "" || s == "-"
}

func StrToUint32(s string) (r uint32) {
	a, _ := strconv.ParseUint(s, 10, 32)
	r = uint32(a)
	return
}
func StrToInt32(s string) (r int32) {
	a, _ := strconv.ParseUint(s, 10, 32)
	r = int32(a)
	return
}

func StrToUint64(s string) (r uint64) {
	r, _ = strconv.ParseUint(s, 10, 64)
	return
}

func StrFromUint32(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

func TrimSpaceArr(s []string) (out []string) {
	out = make([]string, len(s))
	for i, v := range s {
		out[i] = strings.TrimSpace(v)
	}
	return
}

func StrHasLowerPrefix(s string) bool {
	return strings.Contains("abcdefghijklmnopqrstuvwxyz", string(s[0]))
}
