package ucore

import "time"

func NowPtr() *time.Time {
	t := time.Now()
	return &t
}
func TimeToPtr(t time.Time) *time.Time {
	return &t
}
