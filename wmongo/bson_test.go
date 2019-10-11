package wmongo

import (
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	mbson "go.mongodb.org/mongo-driver/bson"
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
	bt, err = bson.MarshalJSON(bm)
	assert.NoError(t, err)
	t.Logf("%s", bt)
}
