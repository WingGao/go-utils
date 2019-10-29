package wmongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"testing"
	"time"
)

type ModelMap struct {
	//MgModel      `bson:",inline"`
	FieldStruct MgTimeModel
	//FieldStr     string
	FieldStructArr []EmModel
	FieldStrB      string `bson:"field_b"`
}

func TestNewMapValueWriter(t *testing.T) {
	mvw := NewMapValueWriter()
	ec := bsoncodec.EncodeContext{Registry: bson.DefaultRegistry}
	enc := &bson.Encoder{}
	enc.SetContext(ec)
	enc.Reset(mvw)
	now := time.Now()
	mod := &ModelMap{
		//FieldStr: "1",
		FieldStruct:    MgTimeModel{&now},
		FieldStrB:      "b",
		FieldStructArr: []EmModel{{FieldA: "c"}, {FieldA: "d"}}}
	err := enc.Encode(mod)
	assert.NoError(t, err)
	t.Log(mvw.buf)
}
