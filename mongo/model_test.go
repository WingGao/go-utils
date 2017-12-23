package mongo

import (
	"os"
	"testing"
	"github.com/WingGao/go-utils"
	"gopkg.in/mgo.v2"
	"github.com/stretchr/testify/assert"
	"mtest"
	"github.com/jinzhu/copier"
)

var (
	testConfig utils.MConfig
	msess      *mgo.Session
)

func TestMain(m *testing.M) {
	testConfig, _ = utils.LoadConfig(os.Getenv("WING_GO_CONF"))
	mgo.SetDebug(true)
	mgo.SetLogger(GetLogger())
	msess, _ = mgo.Dial(testConfig.Mongodb)
	os.Exit(m.Run())
}

type ModelA struct {
	MgModel     `bson:",inline"`
	MgTimeModel `bson:",inline"`
	FieldStr     string
	FieldStructs []EmModel
}
type EmModel struct {
	FieldA string
}

func (ModelA) TableName() string {
	return "test_modela"
}

func NewModelA() (m *ModelA) {
	m = &ModelA{MgModel: MgModel{Session: msess, DbName: "nxpt_dev"}}
	m.SetParent(m)
	return
}
func testP(p IMgParent) {
	p.FormatError(nil)
}
func TestMgModel_Save(t *testing.T) {
	mod := NewModelA()
	testP(mod)
	mod.FieldStr = "1"
	mod.FieldStructs = []EmModel{EmModel{FieldA: "11"}, EmModel{FieldA: "22"},}
	err := mod.Save()
	assert.NoError(t, err)
	assert.NotEmpty(t, mod.Id)
	id1 := mod.Id
	mod.FieldStr = "2"
	err = mod.Save()
	assert.NoError(t, err)
	assert.Equal(t, id1, mod.Id)
	mtest.OutputJson(mod)
}

func TestMgModel_LoadById(t *testing.T) {
	mod := NewModelA()
	err := mod.LoadById("59fae3686154b3790cdc7f81")
	assert.NoError(t, err)
	assert.Equal(t, "2", mod.FieldStr)
	mtest.OutputJson(mod)
}

func TestToObjectId(t *testing.T) {
	mod := NewModelA()
	a := struct {
		Id string
	}{"5a16b6117aff15ead7d498da"}
	copier.Copy(mod, a)
	mod.Id = ToObjectId(mod.Id)
	assert.Equal(t, "5a16b6117aff15ead7d498da", mod.Id.Hex())
}

func TestGetMSetIgnore(t *testing.T) {
	mod := NewModelA()
	mod.AutoNow()
	upm := GetMSetIgnore(mod,"fieldstr")
	t.Log(MarshalJSONStr(upm))
}
