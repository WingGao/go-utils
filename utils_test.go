package utils

import (
	"testing"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"path/filepath"
	"fmt"
	"runtime"
)

var (
	testConfig MConfig
	_db        *gorm.DB
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	host, _ := os.Hostname()
	testEnvFile, _ := filepath.Abs(filepath.Join(filename, fmt.Sprintf("../_tests/config_test.%s.yaml", host)))
	testConfig, _ = LoadConfig(testEnvFile)
	db, err := gorm.Open("mysql", testConfig.GetMySQLString(""))
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	_db = db
	os.Exit(m.Run())
}
