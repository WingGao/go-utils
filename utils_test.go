package utils

import (
	"testing"
	"os"
)

var (
	testConfig MConfig
)

func TestMain(m *testing.M) {
	testConfig, _ = LoadConfig(os.Getenv("WING_GO_CONF"))
	os.Exit(m.Run())
}
