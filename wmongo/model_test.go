package wmongo

import (
	"context"
	"github.com/WingGao/go-utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var (
	testConfig utils.MConfig
	msess      *mongo.Client
	lastId     primitive.ObjectID
)

func TestMain(m *testing.M) {
	msess, _ = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	os.Exit(m.Run())
}

type ModelA struct {
	MgModel      `bson:",inline"`
	MgTimeModel  `bson:",inline"`
	FieldStr     string //field_str
	FieldStructs []EmModel //field_structs
	FieldStrB    string `bson:"field_b"`
	FieldStrList []string
}
type EmModel struct {
	FieldA string
}

func (ModelA) TableName() string {
	return "test_modela"
}

func NewModelA() (m *ModelA) {
	m = &ModelA{MgModel: MgModel{Client: msess, DbName: "go_test"}}
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
	mod.FieldStrB = "alias"
	err := mod.Save()
	assert.NoError(t, err)
	assert.NotEmpty(t, mod.Id)
	assert.NotEmpty(t, mod.UpdatedAt)
	id1 := mod.Id
	mod.FieldStr = "2"
	err = mod.Save()
	assert.NoError(t, err)
	assert.Equal(t, id1, mod.Id)
	lastId = id1
}

func TestMgModel_LoadById(t *testing.T) {
	mod := NewModelA()
	err := mod.LoadById(lastId)
	assert.NoError(t, err)
	assert.Equal(t, "2", mod.FieldStr)
}

func TestToObjectId(t *testing.T) {
	idHex := "5ccad2fb3825f66ccd642ebe"
	//id, _ := primitive.ObjectIDFromHex(idHex)
	assert.Equal(t, 24, len(idHex))
	id2 := ToObjectId(idHex)
	assert.Equal(t, idHex, id2.Hex())
	assert.Equal(t, idHex, ToObjectId(id2).Hex())
}

func TestGetMSetIgnore(t *testing.T) {
	mod := NewModelA()
	mod.AutoNow()
	mod.FieldStr = "123"
	mod.FieldStructs = []EmModel{{"1"}, {"2"}}
	upm, err1 := GetMSetIgnore(mod, "field_str")
	assert.NoError(t, err1)
	assert.NotContains(t, upm["$set"], "field_str")
	assert.Contains(t, upm["$set"], "UpdatedAt")
	t.Log(upm)
	upm, err1 = GetMSetIgnore(mod, " ")
	assert.NoError(t, err1)
	assert.Contains(t, upm["$set"], "field_str")

}
