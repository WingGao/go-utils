package internal

import (
	"fmt"
	"github.com/WingGao/go-utils"
	"github.com/WingGao/go-utils/redis"
	_ "github.com/WingGao/go-utils/wlog"
	"github.com/jinzhu/gorm"
	"github.com/ungerik/go-dry"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"runtime"
)

var (
	TestConfig utils.MConfig
	TestDb     *gorm.DB
)

func LoadFromFile() {
	_, filename, _, _ := runtime.Caller(0)
	host, _ := os.Hostname()
	testEnvFile, _ := filepath.Abs(filepath.Join(filename, fmt.Sprintf("../../_tests/config_test.yaml")))
	tf, _ := dry.FileGetBytes(testEnvFile)
	cm := make(map[interface{}]interface{})
	yaml.Unmarshal(tf, &cm)
	envFilePath := cm[host]
	TestConfig, _ = utils.LoadConfig(envFilePath.(string))
}

func LoadRedis(){
	redis.LoadClient(TestConfig.Redis)
}