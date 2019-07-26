package ucore

import "github.com/json-iterator/go"

func JsonMarshalToString(obj interface{}) string {
	out, err := jsoniter.MarshalToString(obj)
	if err != nil {
		return err.Error()
	}
	return out
}

func JsonMarshalIndentString(obj interface{}, prefix, indent string) string {
	out, err := jsoniter.MarshalIndent(obj, prefix, indent)
	if err != nil {
		return err.Error()
	}
	return string(out)
}
