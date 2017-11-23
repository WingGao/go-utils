package utils

import (
	"testing"
	"github.com/ungerik/go-dry"
	"github.com/stretchr/testify/assert"
)

func TestMConfig_Get(t *testing.T) {
	conf, err := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	assert.NoError(t, err)
	islocal := conf.Get("islocal")
	assert.Equal(t, true, islocal)
}

func TestMConfig_AbsPath(t *testing.T) {
	conf, err := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	assert.NoError(t, err)
	t.Log(conf.configPath)
	t.Log(conf.AbsPath("../helloAbs"))
}
