package utils

import (
	"github.com/thoas/go-funk"
	"github.com/emirpasic/gods/maps/hashmap"
	"reflect"
)

func DefaultUint32(v, def uint32) uint32 {
	if v <= 0 {
		return def
	}
	return v
}

func DefaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func ContainsUint32(s []uint32, v uint32) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}

func MapGetColumn(arr interface{}, key string) interface{} {
	return funk.Map(arr, func(item interface{}) interface{} {
		return funk.Get(item, key)
	})
}
func MapGetColumnUint32(arr interface{}, key string) []uint32 {
	return funk.Map(arr, func(item interface{}) uint32 {
		return funk.Get(item, key).(uint32)
	}).([]uint32)
}
func MapGetColumnString(arr interface{}, key string) []string {
	return funk.Map(arr, func(item interface{}) string {
		return funk.Get(item, key).(string)
	}).([]string)
}

func ArrayToHashmap(arr interface{}, mapKeyField string) *hashmap.Map {
	m := hashmap.New()
	funk.ForEach(arr, func(v interface{}) {
		m.Put(funk.Get(v, mapKeyField), v)
	})
	return m
}

func EqualValUint32(a *uint32, b uint32) bool {
	return a != nil && *a == b
}

//将zero转为nil
func toPrtZero(ptr interface{}) interface{} {
	val := reflect.ValueOf(ptr)
	if val.IsNil() {
		return nil
	}
	if funk.IsZero(val.Elem().Interface()) {
		return nil
	}
	return ptr
}

func ToPrtZeroUint32(ptr uint32) *uint32 {
	if ptr == 0 {
		return nil
	}
	n := ptr
	return &n
}