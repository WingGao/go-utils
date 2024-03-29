package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/rs/xid"
)

// 其他域做主键的测试
type ModPrimaryKey1 struct {
	Model           `gorm:"-"`
	Serial   string `gorm:"primary_key"`
	//Index2   int    `gorm:"primary_key;auto_increment:false"`
	IntField int
}

func (ModPrimaryKey1) TableName() string {
	return "test_mod_primary_key1"
}

func (m *ModPrimaryKey1) SetPrimaryKey() (interface{}, error) {
	if m.Serial == "" {
		m.Serial = xid.New().String()
	}
	return m.Serial, nil
}

func getOneModPrimaryKey1() *ModPrimaryKey1 {
	m := NewModPrimaryKey1()
	ml, _ := m.FindList()
	mod2 := (*ml.(*[]*ModPrimaryKey1))[0]
	return mod2
}

func NewModPrimaryKey1() (m *ModPrimaryKey1) {
	m = &ModPrimaryKey1{}
	m.SetDB(_db)
	m.SetParent(m)
	return
}

func TestAutoMerge(t *testing.T) {
	m := &ModPrimaryKey1{}
	_db.DropTableIfExists(m)
	err := _db.AutoMigrate(m).Error
	assert.NoError(t, err)
}

func TestModel_Save(t *testing.T) {
	mod := NewModPrimaryKey1()
	mod.IntField = 12
	err := mod.UpsertLight()
	assert.NoError(t, err)
	//根系
	mod.IntField = 13
	err = mod.UpsertLight()
	assert.NoError(t, err)
}

func TestModel_LoadByPk(t *testing.T) {
	mod := NewModPrimaryKey1()
	mod2 := getOneModPrimaryKey1()
	err := mod.LoadByPk(mod2.Serial)
	assert.NoError(t, err)
}

func TestModel_Delete(t *testing.T) {
	mod := NewModPrimaryKey1()
	mod.UpsertLight()
	pk := mod.Serial
	mod2 := NewModPrimaryKey1()
	mod2.Serial = pk
	err := mod2.Delete()
	assert.NoError(t, err)
	mod3 := NewModPrimaryKey1()
	err = mod3.LoadByPk(pk)
	assert.Error(t, err)
	mod.Serial = ""
	err = mod.Delete()
	assert.Error(t, err)
}
