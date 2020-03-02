package ucore

import "time"

func NowPtr() *time.Time {
	t := time.Now()
	return &t
}
