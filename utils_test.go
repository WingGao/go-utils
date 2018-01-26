package utils

import (
	"testing"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	testConfig MConfig
	_db        *gorm.DB
)

func TestMain(m *testing.M) {
	testConfig, _ = LoadConfig(os.Getenv("WING_GO_CONF"))
	db, err := gorm.Open("mysql", testConfig.Mysql)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	_db = db
	os.Exit(m.Run())
}
