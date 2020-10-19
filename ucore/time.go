package ucore

import "time"

func NowPtr() *time.Time {
	t := time.Now()
	return &t
}
func TimeToPtr(t time.Time) *time.Time {
	return &t
}

//从t1到现在，必须大于dur
func WaitBefore(t1 time.Time, dur time.Duration) {
	delay := time.Now().Sub(t1)
	if delay < dur { //需要等待
		time.Sleep(delay - dur)
	}
}
