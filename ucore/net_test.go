package ucore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRealIP(t *testing.T) {
	ip, err := GetRealIP()
	assert.NoError(t, err)
	t.Log(ip)
}
