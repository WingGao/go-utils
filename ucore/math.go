package ucore

import (
	"math"
	"time"
)

func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}
func MinInt32(a, b int32) int32 {
	return int32(math.Min(float64(a), float64(b)))
}
func MinInt64(a, b int64) int64 {
	return int64(math.Min(float64(a), float64(b)))
}
func MinUint32(a, b uint32) uint32 {
	return uint32(math.Min(float64(a), float64(b)))
}
func MinUint64(a, b uint64) uint64 {
	return uint64(math.Min(float64(a), float64(b)))
}
func MaxUint32(a, b uint32) uint32 {
	return uint32(math.Max(float64(a), float64(b)))
}
func MaxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
func MaxInt64(a, b int64) int64 {
	return int64(math.Max(float64(a), float64(b)))
}

// 获取周岁
func CalcAge(birthday time.Time) uint32 {
	now := time.Now()

	if birthday.After(now) { //溢出
		return 0
	}
	// 周岁＝今年-出生年（已过生日）（未过生日还要-1）
	year := now.Year() - birthday.Year()
	if now.Month() > birthday.Month() || (now.Month() == birthday.Month() && now.Day() >= birthday.Day()) {
		//已过生日
	} else {
		year -= 1
	}
	return uint32(year)
}
