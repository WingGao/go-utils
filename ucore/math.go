package ucore

import (
	"math"
	"time"
)

func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

// 获取周岁
func CalcAge(birthday time.Time) uint32 {
	now := time.Now()
	// 周岁＝今年-出生年（已过生日）（未过生日还要-1）
	year := now.Year() - birthday.Year()
	if now.Month() >= birthday.Month() && now.Day() >= birthday.Day() {
		//已过生日
	} else {
		year -= 1
	}
	return uint32(year)
}
