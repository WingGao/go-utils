package utils

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"encoding/json"
	"github.com/go-playground/form"
	"github.com/json-iterator/go"
	"github.com/emirpasic/gods/sets/hashset"
	"net/url"
	"net/http"
)

var ignoreErros = hashset.New()

const (
	HANDLER_CANCEL = "wing-handler-cancel"
)

type ErrJson struct {
	Err string `json:"err_msg"`
}

func AddHandlerIgnoreErrors(errs ...interface{}) {
	ignoreErros.Add(errs...)
}
func CancelAfterHandler(ictx context.Context) {
	ictx.Values().Set(HANDLER_CANCEL, true)
}

func AfterHandler(ictx context.Context, o interface{}, err error) {
	FixSetCookie(ictx)
	//跳过已经处理过的请求
	if v := ictx.Values().Get(HANDLER_CANCEL); v != nil && v.(bool) {
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
		err = NewErrParams()
	}
	return
}

//params必须是ptr
func ParseForm(v url.Values, outPtr interface{}) (err error) {
	dec := form.NewDecoder()
	err = dec.Decode(outPtr, v)
	return
}

// 只保留最新的Set-Cookie
func FixSetCookie(ctx context.Context) {
	respHeader := ctx.ResponseWriter().Header()
	hresp := http.Response{Header: respHeader}
	cookies := hresp.Cookies()
	cookieNum := len(cookies)
	if cookieNum >= 2 {
		cookieMap := make(map[string]*http.Cookie, cookieNum)
		for _, cookie := range cookies {
			cookieMap[cookie.Name] = cookie
		}
		//有重复
		if len(cookieMap) < cookieNum {
			ctx.Header("Set-Cookie", "")
			for _, cookie := range cookieMap {
				ctx.SetCookie(cookie)
			}
		}
	}
}

func FormHasValue(ctx context.Context, name string) bool {
	fv := ctx.FormValues()
	if fv == nil {
		return false
	}
	if v, ok := fv[name]; ok && len(v) > 0 {
		return true
	}
	return false
}
