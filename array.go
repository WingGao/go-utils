package utils

import (
	"github.com/thoas/go-funk"
	"reflect"
)

func GetItemString(arr []string, index int) (out string, ok bool) {
	if len(arr) > index {
		out = arr[index]
		ok = true
	}
	return
}

//
func ArrJoin(arrL interface{}, arrR interface{}, keyL, keyR string, joinFunc interface{}) interface{} {
	rmap := ArrayToHashmap(arrR, keyR)
	funcVal := reflect.ValueOf(joinFunc)
	funk.ForEach(arrL, func(vl interface{}) {
		keyV := funk.Get(vl, keyL)
		if vr, ok := rmap.Get(keyV); ok {
			funcVal.Call([]reflect.Value{reflect.ValueOf(vl), reflect.ValueOf(vr)})
		}
	})
	return arrL
}
