package wmongo

import (
	mbson "github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBSet(t *testing.T) {
	modA := NewModelA()
	modA.FieldStr = "123"
	modA.Id = primitive.NewObjectID()
	bm := BSet(modA)
	bt, err := bson.Marshal(bm)
	assert.NoError(t, err)
	t.Logf("%s", bt)

	bt2, err2 := mbson.Marshal(bm)
	assert.NoError(t, err2)
	t.Logf("%s", bt2)
	bt, err = bson.MarshalExtJSON(bm, true, false)
	assert.NoError(t, err)
	t.Logf("%s", bt)
}

func TestIgnore(t *testing.T) {
	modA := NewModelA()
	modA.FieldStr = "123"
	modA.FieldStrB = "234"
	v, err := GetMSetIgnore(modA, "field_b")
	assert.NoError(t, err)
	t.Logf("%v", v)
	v, err = GetBsonM(modA, "field_b")
	assert.NoError(t, err)
	t.Logf("%v", v)
}

func TestStructToBsonMap(t *testing.T) {
	modA := NewModelA()
	modA.FieldStr = "123"
	modA.FieldStrList = []string{"Ph0jiw5h", "UiGGQ3Y-5", "Xtpe0cyZd", "rDARrNQrD", "zGRBl01PK", "0INuPBdxh", "E1vHjXJkT", "gAEUVu0wH"}
	v, err := GetBsonM(modA)
	assert.NoError(t, err)
	t.Logf("%v", v)
}

