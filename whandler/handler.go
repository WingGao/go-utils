package whandler

import (
	"encoding/json"
	ucore "github.com/WingGao/go-utils/ucore"
	"github.com/WingGao/go-utils/wlog"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-playground/form"
	"github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"net/http"
	"net/url"
	"strings"
)

var ignoreErros = hashset.New()

const (
	HANDLER_CANCEL = "wing-handler-cancel"
)

// 自己处理json，不适用默认的处理
type JsonRep struct {
	Json []byte
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
		var returnErr interface{}
		if !ignoreErros.Contains(err.Error()) {
			err2 := ucore.NewWError(err)
			err2.Fmt()
			wlog.S().Error(err2.Err.Error(), "\n", err2.ErrorStack())
		}
		ictx.StatusCode(iris.StatusBadRequest)
		//if err3, ok := err.(*errors.Error); ok {
		//	if err3, ok1 := err3.Err.(ucore.CommError); ok1 {
		//		returnErr = err3
		//	}
		//}
		if err3, ok1 := err.(ucore.CommError); ok1 {
			returnErr = err3
		}
		if returnErr == nil {
			returnErr = ucore.CommError{ErrMsg: err.Error()}
		}
		ictx.JSON(returnErr)
	} else {
		var buf []byte

		if offj, isffj := o.(json.Marshaler); isffj {
			buf, err = offj.MarshalJSON()
		} else if jp, ok := o.(JsonRep); ok {
			buf = jp.Json
		} else {
			buf, err = jsoniter.Marshal(o)
		}

		if err != nil {
			ictx.StatusCode(iris.StatusBadRequest)
			ictx.JSON(ucore.CommError{ErrMsg: err.Error()})
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
		err = ucore.NewErrParams()
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

// 替换路由
func ReplaceRoute(app *iris.Application, r *router.Route) {
	old := app.GetRoute(r.Name)
	if old.MainHandlerName != r.MainHandlerName {
		//替换
		routers := app.GetRoutes()
		for i, v := range routers {
			if v.Name == r.Name {
				routers[i] = r
			}
		}
	}
}

// 合并同路由
func RouteAddHandlers(app *iris.Application, method, subdomain, unparsedPath string, handlers ...context.Handler) {
	old := app.GetRoute(method + subdomain + unparsedPath)
	old.Handlers = append(old.Handlers, handlers...) // 添加最后1个
}

func GetHandlerIp(c context.Context) string {
	if ip := c.Request().Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 && ips[0] != "" {
			rip := strings.Split(ips[0], ":")
			return rip[0]
		}
	} else {
		ip := strings.Split(c.Request().RemoteAddr, ":")
		if len(ip) > 0 {
			if ip[0] != "[" {
				return ip[0]
			}
		}
	}
	return ""
}
