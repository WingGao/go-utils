package utils

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"encoding/json"
	"github.com/go-playground/form"
	"github.com/json-iterator/go"
)

type ErrJson struct {
	Err string `json:"err_msg"`
}

func AfterHandler(ictx context.Context, o interface{}, err error) {
	if err != nil {
		ictx.StatusCode(iris.StatusBadRequest)
		ictx.JSON(ErrJson{Err: err.Error()})
	} else {
		var buf []byte
		offj, isffj := o.(json.Marshaler)
		if isffj {
			buf, err = offj.MarshalJSON()
		} else {
			buf, err = jsoniter.Marshal(o)
		}

		if err != nil {
			ictx.StatusCode(iris.StatusBadRequest)
			ictx.JSON(ErrJson{Err: err.Error()})
		} else {
			ictx.StatusCode(iris.StatusOK)
			ictx.ContentType("application/json")
			ictx.Write(buf)
		}
	}
	//ictx.Next()
}

//params必须是ptr
func ParseForm(ictx context.Context, params interface{}) (err error) {
	dec := form.NewDecoder()
	err = dec.Decode(params, ictx.FormValues())
	if err != nil {
		ictx.Application().Logger().Warnf("Error when reading form: " + err.Error())
		err = ERR_PARAMS
	}
	return
}
