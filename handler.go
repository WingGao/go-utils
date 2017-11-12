package utils

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"encoding/json"
	"github.com/go-playground/form"
	"github.com/json-iterator/go"
	"github.com/emirpasic/gods/sets/hashset"
	"net/url"
)

var ignoreErros = hashset.New()

type ErrJson struct {
	Err string `json:"err_msg"`
}

func AddHandlerIgnoreErrors(errs ...interface{}) {
	ignoreErros.Add(errs...)
}

func AfterHandler(ictx context.Context, o interface{}, err error) {
	//跳过已经处理过的请求
	if ictx.GetStatusCode() > 200 {
		return
	}
	if err != nil {
		if !ignoreErros.Contains(err) {
			err2 := NewWError(err)
			err2.Fmt()
			ictx.Application().Logger().Error(err2.ErrorStack())
		}
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
func ParseFormIris(ictx context.Context, params interface{}) (err error) {
	dec := form.NewDecoder()
	err = dec.Decode(params, ictx.FormValues())
	if err != nil {
		ictx.Application().Logger().Warnf("Error when reading form: " + err.Error())
		err = ERR_PARAMS
	}
	return
}

//params必须是ptr
func ParseForm(v url.Values, outPtr interface{}) (err error) {
	dec := form.NewDecoder()
	err = dec.Decode(outPtr, v)
	return
}
