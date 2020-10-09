package utils

import (
	"runtime"
	"testing"
	"github.com/ungerik/go-dry"
	"github.com/stretchr/testify/assert"
	"net"
	"fmt"
)

func TestMConfig_Get(t *testing.T) {
	conf, err := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	assert.NoError(t, err)
	islocal := conf.Get("islocal")
	assert.Equal(t, true, islocal)
}

func TestMConfig_GetString(t *testing.T) {
	conf, err := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	assert.NoError(t, err)
	s := conf.GetString("test.a", "")
	assert.Equal(t, "a123", s)
	s = conf.GetString("testadminsession", "")
	assert.Equal(t, "d7e3c712-15ad-4935-a7c9-6d4fc1719fee", s)
}

func TestMConfig_AbsPath(t *testing.T) {
	conf, err := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	assert.NoError(t, err)
	t.Log(conf.configPath)
	t.Log(conf.AbsPath("../helloAbs"))
}

func TestMConfig_Addr(t *testing.T) {
	conf, _ := NewConfigFromFile(dry.GetenvDefault("WING_GO_CONF", ""))
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port))
	assert.Equal(t, 7031, addr.Port)
}

func TestLoadForTest(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
}

