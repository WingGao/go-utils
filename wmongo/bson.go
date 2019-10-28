package wmongo

import (
	"github.com/WingGao/go-utils"
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
func BOr(items ...bson.M) (bson.M) {
	return bson.M{"$or": items}
}
func BAnd(items ...bson.M) (bson.M) {
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
func BAvg(v interface{}) (bson.M) {
	return bson.M{"$avg": v}
}

func BAddFields(field string, v interface{}) (out bson.M) {
	return bson.M{"$addFields": bson.M{field: v}}
}

func BMatch(v interface{}) (bson.M) {
	return bson.M{"$match": v}
}
func BGroup(v interface{}) (bson.M) {
	return bson.M{"$group": v}
}

func BIn(v interface{}) (out bson.M) {
	return bson.M{"$in": v}
}

func BInField(field string, v interface{}) (out bson.M) {
	return bson.M{field: bson.M{"$in": v}}
}

func BExists(v interface{}) (bson.M) {
	return bson.M{"$exists": v}
}

//array
func BElemMatch(v interface{}) (out bson.M) {
	return bson.M{"$elemMatch": v}
}

//忽略某些
func GetMSetIgnore(obj interface{}, bsonFields ...string) (bm bson.M) {
	//setM := bson.M{}
	//objt := reflect.TypeOf(obj)
	//objv := reflect.ValueOf(obj)
	//if objt.Kind() == reflect.Ptr {
	//	objt = objt.Elem()
	//	objv = objv.Elem()
	//}
	//if err1 != nil {
	//	return
	//}
	//ignoreMap := make(map[string]bool, len(bsonFields))
	//for _, f := range bsonFields {
	//	ignoreMap[f] = true
	//}
	//
	//for _, v := range info.FieldsList {
	//	if v.Key == "_id" { //忽略_id
	//		continue
	//	}
	//	if _, ok := ignoreMap[v.Key]; !ok {
	//		setv := objv
	//		if len(v.Inline) > 0 {
	//			//inline
	//			for _, inlineNum := range v.Inline {
	//				setv = setv.Field(inlineNum)
	//			}
	//		} else {
	//			setv = setv.Field(v.Num)
	//		}
	//		setM[v.Key] = setv.Interface()
	//	}
	//}
	////panic("not implement")
	//bm = BSet(setM)

	return
}

func BMarshal(m interface{}) []byte {
	bs, _ := bson.Marshal(m)
	return bs
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
