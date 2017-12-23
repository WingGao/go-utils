package mongo

import "gopkg.in/mgo.v2/bson"

func MarshalJSONStr(m interface{}) string {
	js, _ := bson.MarshalJSON(m)
	return string(js)
}
