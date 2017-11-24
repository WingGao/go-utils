package utils

import (
	"testing"
	"reflect"
	"time"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

type ObjInner struct {
	ID int
}
type ObjInner2 struct {
	ID2 int
}
type _obj1 struct {
	ObjInner
	ID2       int
	Attr1     string
	Attr2     int
	Attr3     bool
	Attrb     int
	AttrTime1 time.Time
	AttrTime2 int64
}

type _obj2 struct {
	ObjInner2
	ID        int
	Attr1     string
	Attr2     int
	Attr3     bool
	Attra     string
	AttrTime1 int64
	AttrTime2 time.Time
}

type structInner1 struct {
	ObjInner
	ID2   int
	Attr1 string
	Attr2 int
	Attr3 bool
	Attrb int
}
type structInner2 struct {
	ObjInner2
	ID    int
	Attr1 string
	Attr2 int
	Attr3 bool
	Attra string
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
}
type structSameName2 struct {
	A string
	B time.Time
	C int64
}

func printFields(t *testing.T, val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		tfield := val.Type().Field(i)
		fmt.Sprintf("a", "b")
		t.Logf("%s %s", tfield.Name, tfield.Type)
	}
}

type objectField struct {
	field     reflect.Value
	fieldType reflect.StructField
	parent    reflect.StructField
}

func DeepFields(iface interface{}, parent reflect.StructField) (map[string]objectField) {
	objectFields := make(map[string]objectField, 10)

	ifv := reflect.ValueOf(iface)
	if ifv.Kind() == reflect.Ptr {
		ifv = reflect.Indirect(ifv)
	}
	ift := ifv.Type()
	startIndex := 0
	if ifv.Field(0).Kind() == reflect.Struct && ift.Field(0).Anonymous {
		nfields := DeepFields(ifv.Field(0).Interface(), ift.Field(0))
		for k, v := range nfields {
			objectFields[k] = v
		}
		startIndex = 1
	}
	for i := startIndex; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		ft := ift.Field(i)
		k := ft.Name
		field := objectField{field: v, fieldType: ft, parent: parent}
		objectFields[k] = field
	}

	return objectFields
}

func TestDeepFields(t *testing.T) {
	obj2 := &_obj2{}
	tarFields := DeepFields(obj2, reflect.StructField{})
	t.Logf("%#v", tarFields)
}

//func TestCopyAttrs(t *testing.T) {
//	obj1 := _obj1{
//		ObjInner:  ObjInner{ID: 12},
//		ID2:       13,
//		Attr1:     "a",
//		Attr2:     2,
//		Attr3:     true,
//		Attrb:     4,
//		AttrTime1: time.Now(),
//	}
//	obj2 := &_obj2{}
//
//	//tarValue := reflect.ValueOf(obj2)
//	//tarValue = reflect.Indirect(tarValue)
//	//srcValue := reflect.Indirect(reflect.ValueOf(obj1))
//	//printFields(t, tarValue)
//
//	srcFields, _ := DeepFields(obj1, reflect.StructField{})
//	tarFields, tarTypes := DeepFields(obj2, reflect.StructField{})
//	t.Logf("src=> %#v", srcFields)
//	t.Logf("tar=> %#v", tarFields)
//	for k, v := range srcFields {
//		tarFieldVal := tarFields[k]
//		tType, ok := tarTypes[k]
//		sfval := v
//		if sfval.IsValid() && ok {
//			if sfval.Type() == tarFieldVal.Type() {
//				if !tarFieldVal.CanSet() {
//					//tarFieldVal = tarFieldVal.Addr()
//					t.Logf("cann't set field %s\n%#v\n%#v", tarTypes[k].Name, tarFieldVal, tType)
//				}
//				tarFieldVal.Set(sfval)
//			} else if tType.Type.Name() == "Time" && sfval.Type().Name() == "int64" {
//				// int64 => time.Time
//				t := time.Unix(sfval.Int(), 0)
//				tarFieldVal.Set(reflect.ValueOf(t))
//			} else if tType.Type.Name() == "int64" && sfval.Type().Name() == "Time" {
//				// time.Time => int64
//				tarFieldVal.SetInt(sfval.Interface().(time.Time).Unix())
//			} else {
//				t.Logf("%s=%#v", tType.Name, tType.Type.Name())
//			}
//		}
//	}
//	//t.Log(tarValue.NumField())
//	//CopyAttrs(obj1, obj2)
//	t.Logf("%#v", obj2)
//}

func TestMergo(t *testing.T) {
	obj1 := _obj1{
		ObjInner:  ObjInner{ID: 12},
		ID2:       13,
		Attr1:     "a",
		Attr2:     2,
		Attr3:     true,
		Attrb:     4,
		AttrTime1: time.Now(),
	}
	obj2 := &_obj2{}
	err := mergo.MergeWithOverwrite(obj2, obj1)
	t.Log(err)
	t.Logf("%#v", obj2)
}

func TestCopier(t *testing.T) {
	//obj1 := structInner1{
	//	ObjInner: ObjInner{ID: 12},
	//	ID2:      13,
	//	Attr1:    "a",
	//	Attr2:    2,
	//	Attr3:    true,
	//	Attrb:    4,
	//}
	//obj2 := &structInner2{}
	//err := copier.Copy(obj2, &obj1)
	//t.Log(err)
	//t.Logf("%#v", obj2)

	objs1 := structSameName1{A: "123", B: 2, C: time.Now()}
	objs2 := &structSameName2{}
	err := copier.Copy(objs2, &objs1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", objs2)
	objs3 := make(map[string]interface{})
	err = copier.Copy(&objs3, objs1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", objs3)
}

func TestToPrtZero(t *testing.T) {
	var aUint32 uint32 = 0
	res := ToPrtZeroUint32(aUint32)
	assert.Nil(t, res)
	var bUint32 uint32 = 1
	res = ToPrtZeroUint32(bUint32)
	assert.EqualValues(t, bUint32, *res)
}

func TestObjectGet(t *testing.T) {
	r := ObjectGet("a", func(v string) bool {
		return v == "a"
	})
	assert.Equal(t, true, r)
	t.Log(r)
}
