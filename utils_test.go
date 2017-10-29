package utils

import (
	"testing"
	"os"
)

var (
	testConfig MConfig
)

func TestMain(m *testing.M) {
	//TODO 修改测试配置
	testConfig, _ = LoadConfig("/Users/suamo/Projs/SuamoIris/dev/config.imac.yaml")
	os.Exit(m.Run())
}
