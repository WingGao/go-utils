package ucore

import (
	"math/rand"
	"time"
	"strconv"
	"strings"
	"github.com/thoas/go-funk"
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
	sr := []rune(s) // for unicode
	length = MinInt(length, len(sr))
	if length == 0 {
		return ""
	}

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

//数组
func StrToUint32Arr(s string) (r []uint32) {
	return funk.Map(strings.Split(s, ","), StrToUint32).([]uint32)
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

func DesensitizePhone(phone string) string {
	if len(phone) >= 11 {
		sb := &StringBuilder{}
		sb.Write(phone[:3])
		for i := 0; i < len(phone)-3-4; i++ {
			sb.Write("*")
		}
		sb.Write(phone[len(phone)-4:])
		return sb.String()
	}
	return phone
}
