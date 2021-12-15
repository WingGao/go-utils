package ucore

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
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

// 数组转hashmap
// `mapKeyField`可以是`string`或者是`func`
func ArrayToHashmap(arr interface{}, mapKeyField interface{}) *hashmap.Map {
	m := hashmap.New()
	funk.ForEach(arr, func(v interface{}) {
		m.Put(ObjectGet(v, mapKeyField), v)
	})
	return m
}

func ArrayGetColumnSet(arr interface{}, key string) *hashset.Set {
	kset := hashset.New()
	funk.ForEach(arr, func(v interface{}) {
		kset.Add(funk.Get(v, key))
	})
	return kset
}

// 将arrR的值通过joinFunc来赋值给arrL, 忽律不匹配的元素
// joinFunc ==> func(left typeA, right typeB) any
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
//func ArrRemoveByIndex(arr interface{}, idx int) interface{}{
//	arrVal := reflect.ValueOf(arr)
//	reflect.SliceOf(arrVal.Slice())
//}

func InsertStr(s []string, k int, vs ...string) []string {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]string, len(s) + len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}
