package ucore

import "testing"

func TestGetRealIP(t *testing.T) {
	ip := GetRealIP()
	t.Log(ip)
}
