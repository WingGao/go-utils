package utils

import (
	"github.com/thoas/go-funk"
	"reflect"
	"github.com/emirpasic/gods/maps/hashmap"
)

func GetItemString(arr []string, index int) (out string, ok bool) {
	if len(arr) > index {
		out = arr[index]
		ok = true
	}
	return
}

// 数组转hashmap
// `mapKeyField`可以是`string`或者是`func`
func ArrayToHashmap(arr interface{}, mapKeyField interface{}) *hashmap.Map {
	m := hashmap.New()
	funk.ForEach(arr, func(v interface{}) {
		m.Put(ObjectGet(v, mapKeyField), v)
	})
	return m
}

//
func ArrJoin(arrL interface{}, arrR interface{}, keyL, keyR interface{}, joinFunc interface{}) interface{} {
	rmap := ArrayToHashmap(arrR, keyR)
	funcVal := reflect.ValueOf(joinFunc)
	funk.ForEach(arrL, func(vl interface{}) {
		keyV := ObjectGet(vl, keyL)
		if vr, ok := rmap.Get(keyV); ok {
			funcVal.Call([]reflect.Value{reflect.ValueOf(vl), reflect.ValueOf(vr)})
		}
	})
	return arrL
}
