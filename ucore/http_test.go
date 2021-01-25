package ucore

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetPublicIP(t *testing.T) {
	ip, err := GetPublicIP()
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)
	t.Log("ip", ip)
}
