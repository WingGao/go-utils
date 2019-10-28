package wmongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"testing"
)

type ModelMap struct {
	MgModel      `bson:",inline"`
	MgTimeModel  `bson:",inline"`
	FieldStr     string
	FieldStructs []EmModel
	FieldStrB    string `bson:"field_b"`
}

func TestNewMapValueWriter(t *testing.T) {
	mvw := NewMapValueWriter()
	ec := bsoncodec.EncodeContext{Registry: bson.DefaultRegistry}
	enc := &bson.Encoder{}
	enc.SetContext(ec)
	enc.Reset(mvw)
	mod := &ModelMap{FieldStr: "1", FieldStrB: "b",
		FieldStructs: []EmModel{{FieldA: "c"}, {FieldA: "d"}}}
	err := enc.Encode(mod)
	assert.NoError(t, err)
	t.Log(mvw.buf)
}
