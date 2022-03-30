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

// t1 到 t2 之间的差
func TimeDiffMs(t1 time.Time, t2 time.Time) int64 {
	return (int64)(t1.Sub(t2) / time.Millisecond)
}
