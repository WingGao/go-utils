package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/ungerik/go-dry"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	testConfig MConfig
	_db        *gorm.DB
)
func loadFromFile(){
	_, filename, _, _ := runtime.Caller(0)
	host, _ := os.Hostname()
	testEnvFile, _ := filepath.Abs(filepath.Join(filename, fmt.Sprintf("../_tests/config_test.yaml")))
	tf, _ := dry.FileGetBytes(testEnvFile)
	cm := make(map[interface{}]interface{})
	yaml.Unmarshal(tf, &cm)
	envFilePath := cm[host]
	testConfig, _ = LoadConfig(envFilePath.(string))
}
func TestMain(m *testing.M) {
	loadFromFile()
	db, err := gorm.Open("mysql", testConfig.GetMySQLString(""))
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	_db = db
	os.Exit(m.Run())
}
