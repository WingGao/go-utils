package wmongo

import (
	"github.com/WingGao/go-utils"
	jsoniter "github.com/json-iterator/go"

	//"github.com/globalsign/mgo/bson"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"reflect"
	"strings"
)

func init() {
	// 替换为snake模式
	reg := bson.NewRegistryBuilder()
	structCodec, err1 := bsoncodec.NewStructCodec(SnakeStructTagParser)
	utils.PanicIfErr(err1)
	reg.RegisterDefaultEncoder(reflect.Struct, structCodec)
	reg.RegisterDefaultDecoder(reflect.Struct, structCodec)
	bson.DefaultRegistry = reg.Build()
}

func MarshalJSONStr(m interface{}) string {
	//js, _ := bson.MarshalExtJSON(m, false, false)
	//return string(js)
	return ""
}

func BEq(v interface{}) (out bson.M) {
	return bson.M{"$eq": v}
}
func BNe(v interface{}) (out bson.M) {
	return bson.M{"$ne": v}
}

//query logical
func BOr(items ...bson.M) bson.M {
	return bson.M{"$or": items}
}
func BAnd(items ...bson.M) bson.M {
	return bson.M{"$and": items}
}

func BSet(v interface{}) (out bson.M) {
	return bson.M{"$set": v}
}
func BUnset(v interface{}) (out bson.M) {
	return bson.M{"$unset": v}
}
func BCount(v interface{}) (out bson.M) {
	return bson.M{"$count": v}
}
func BSum(v interface{}) (out bson.M) {
	return bson.M{"$sum": v}
}
func BAvg(v interface{}) bson.M {
	return bson.M{"$avg": v}
}

func BAddFields(field string, v interface{}) (out bson.M) {
	return bson.M{"$addFields": bson.M{field: v}}
}

func BMatch(v interface{}) bson.M {
	return bson.M{"$match": v}
}
func BGroup(v interface{}) bson.M {
	return bson.M{"$group": v}
}

func BIn(v interface{}) (out bson.M) {
	return bson.M{"$in": v}
}

func BInField(field string, v interface{}) (out bson.M) {
	return bson.M{field: bson.M{"$in": v}}
}

func BExists(v interface{}) bson.M {
	return bson.M{"$exists": v}
}

//array
func BElemMatch(v interface{}) (out bson.M) {
	return bson.M{"$elemMatch": v}
}

//忽略某些 bsonFields是数据库里的名字，不是struct属性
func GetMSetIgnore(obj interface{}, ignoreBsonFields ...string) (bm bson.M, err error) {
	bmap, err1 := StructToBsonMap(obj)
	if err1 != nil {
		return nil, err1
	}
	for _, f := range ignoreBsonFields {
		delete(bmap, f)
	}
	bm = BSet(bmap)
	return
}

//忽略某些 filterBsonFields 是数据库里的名字，不是struct属性, 不传不过滤
func GetBsonM(obj interface{}, filterBsonFields ...string) (bm bson.M, err error) {
	bmap, err1 := StructToBsonMap(obj)
	if err1 != nil {
		return nil, err1
	}
	fields := bson.M{}
	if len(filterBsonFields) > 0 {
		for _, f := range filterBsonFields {
			fields[f] = bmap[f]
		}
	} else {
		fields = bmap
	}

	return fields, err
}

func BMarshal(m interface{}) []byte {
	bs, _ := bson.Marshal(m)
	return bs
}

// 从json中合并到原对象，作用是处理空值
func mergeFromJson(obj interface{}, bs []byte) error {
	err := jsoniter.Unmarshal(bs, obj)
	return err
}
func mergeFromRequest(obj interface{}, bs []byte) error {
	err := jsoniter.Unmarshal(bs, obj)
	return err
}

// 默认是snake_case模式
var SnakeStructTagParser bsoncodec.StructTagParserFunc = func(sf reflect.StructField) (bsoncodec.StructTags, error) {
	key := strcase.ToSnake(sf.Name)
	tag, ok := sf.Tag.Lookup("bson")
	if !ok && !strings.Contains(string(sf.Tag), ":") && len(sf.Tag) > 0 {
		tag = string(sf.Tag)
	}
	var st bsoncodec.StructTags
	if tag == "-" {
		st.Skip = true
		return st, nil
	}

	for idx, str := range strings.Split(tag, ",") {
		if idx == 0 && str != "" {
			key = str
		}
		switch str {
		case "omitempty":
			st.OmitEmpty = true
		case "minsize":
			st.MinSize = true
		case "truncate":
			st.Truncate = true
		case "inline":
			st.Inline = true
		}
	}

	st.Name = key

	return st, nil
}
