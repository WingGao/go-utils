package utils

import "github.com/json-iterator/go"

func JsonMarshalToString(obj interface{}) string {
	out, err := jsoniter.MarshalToString(obj)
	if err != nil {
		return err.Error()
	}
	return out
}
