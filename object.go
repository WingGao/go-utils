package utils

import (
	"github.com/thoas/go-funk"
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

func MapGetColumn(arr interface{}, key string) []interface{} {
	return funk.Map(arr, func(item interface{}) interface{} {
		return funk.Get(item, key)
	}).([]interface{})
}

//转换*uint32和uint32
func MapGetColumnUint32(arr interface{}, key string) []uint32 {
	return funk.Map(arr, func(item interface{}) uint32 {
		v := funk.Get(item, key)
		if vi, ok := v.(*uint32); ok {
			if vi == nil {
				return 0
			}
			return *vi
		}
		return v.(uint32)
	}).([]uint32)
}

func MapGetColumnString(arr interface{}, key string) []string {
	return funk.Map(arr, func(item interface{}) string {
		return funk.Get(item, key).(string)
	}).([]string)
}

func ArrayFilterNotEmpty(arr interface{}) []interface{} {
	return funk.Filter(arr, funk.NotEmpty).([]interface{})
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

func FromPrtZeroUint32(ptr *uint32) uint32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

//只获取类型指针，值不会拷贝
func PtrOf(ob interface{}) (out interface{}) {
	t := reflect.TypeOf(ob)

	// Avoid double pointers if itf is a pointer
	if t.Kind() == reflect.Ptr {
		cp := reflect.New(t.Elem())
		out = cp.Interface()
	} else {
		cp := reflect.New(t)
		out = cp.Interface()
	}

	//err := copier.Copy(out, ob)
	//if err != nil {
	//	panic(err)
	//}
	return
}

func ObjectGet(v, f interface{}) interface{} {
	if k, ok := f.(string); ok {
		return funk.Get(v, k)
	} else {
		ft := reflect.ValueOf(f)
		if ft.Kind() == reflect.Func {
			rv := ft.Call([]reflect.Value{reflect.ValueOf(v)})[0].Interface()
			return rv
		} else {
			panic("ObjectGet need f")
		}
	}
	return nil
}
