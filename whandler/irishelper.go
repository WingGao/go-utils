package whandler

import (
	"github.com/WingGao/errors"
	"github.com/WingGao/go-utils/wlog"
	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"reflect"
	"strings"
)

type IrisHelper struct {
	SessionType reflect.Type
	formDec *form.Decoder
}

// 将友好请求转换为iris
// 支持函数：
// 	func(*ctx,*any)(out,err)
// 	func(*ctx,*any,sess)(out,err)
func (h *IrisHelper) WrapH(handler interface{})iris.Handler {
	if h.formDec == nil {
		h.formDec = form.NewDecoder()
	}
	handlerVal := reflect.ValueOf(handler)
	handlerType := handlerVal.Type()
	typeReq := handlerType.In(1) // *pb.xxxReq

	return func(ictx *context.Context) {
		var output interface{}
		var (
			err  error
			out  []reflect.Value
			sess interface{}
		)
		defer (func() {
			if err2 := recover(); err2 != nil {
				//err = errors.Wrap(err2, 1)
				err = err2.(error)
			}
			AfterHandler(ictx, output, err)
		})()
		// 构造参数
		reqVal := reflect.New(typeReq.Elem())
		if typeReq.Kind() == reflect.Ptr {
			reqVal.Elem().Set(reflect.Zero(typeReq.Elem()))
		} else {
			reqVal.Elem().Set(reflect.Zero(typeReq))
		}
		reqPtr := reqVal.Interface()
		body, _ := ictx.GetBody()
		if ictx.Method() != "GET" && strings.Contains(ictx.GetContentTypeRequested(), context.ContentJSONHeaderValue) && len(body) > 0 {
			err = jsoniter.Unmarshal(body, reqPtr)
		} else {
			err = h.formDec.Decode(reqPtr,ictx.FormValues())
		}
		if err != nil {
			wlog.S().Error("Error when reading form: ", err.Error())
			err = errors.New("解析错误")
			return
		}

		if handlerType.NumIn() == 3 {
			//TODO 带session
			//sess, err = usess.NewSessionFromIris(ictx, utils.XSESSION_KEY)
			out = handlerVal.Call([]reflect.Value{reflect.ValueOf(ictx), reqVal, reflect.ValueOf(sess)})
		} else { //不带session
			out = handlerVal.Call([]reflect.Value{reflect.ValueOf(ictx), reqVal})
		}
		output = out[0].Interface()
		errVal := out[1]
		if !errVal.IsNil() {
			err = errVal.Interface().(error)
		}
		if sess != nil { //处理session
			//sess.SaveIris(ictx, utils.XSESSION_KEY)
		}
	}
}
