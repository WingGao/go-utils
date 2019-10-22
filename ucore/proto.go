package ucore

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

// 将一个struct转换为.proto文件中的message，方便开发
func StructToProtofile(item interface{}, name string) (pbtxt string, err error) {
	mj, err1 := jsoniter.Marshal(item)
	if err1 != nil {
		return pbtxt, err1
	}
	fieldAny := jsoniter.Get(mj)
	fieldMap := make(map[string]interface{})
	err = jsoniter.Unmarshal(mj, fieldMap)
	fmt.Println(fieldMap)
	sb := StringBuilder{}
	sb.WriteF("message %s {\n", name)
	//idx := 1
	for i, k := range fieldAny.Keys() {
		field := ObjectGet(item, k)
		rt := reflect.TypeOf(field)
	LAB_PTR:
		pType := rt.Name()
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
			goto LAB_PTR
		}
		switch rt.Name() {
		case "Time":
			pType = "int64"
		}
		sb.WriteF("%s %s = %d;\n", pType, k, i+1)
	}
	sb.Write("}\n")
	pbtxt = sb.String()
	return
}
