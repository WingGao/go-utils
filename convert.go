package utils

import "github.com/thoas/go-funk"

func ToInterfaceArrayString(t []string) []interface{} {
	s := make([]interface{}, len(t))
	for i, v := range t {
		s[i] = v
	}
	return s
}

func ToInterfaceArray(m interface{}) []interface{} {
	return funk.Map(m, func(x interface{}) interface{} {
		return x
	}).([]interface{})
}
