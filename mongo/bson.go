package mongo

import "github.com/globalsign/mgo/bson"

func MarshalJSONStr(m interface{}) string {
	js, _ := bson.MarshalJSON(m)
	return string(js)
}
