package ucore

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
)

func GetRealIP() (nip net.IP, err error) {
	method := "api"
	ip := ""
	switch method {
	case "api":
		_, body, errs := gorequest.New().Get("http://ip-api.com/json/").End()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		v := make(map[string]interface{}, 20)
		if err = jsoniter.UnmarshalFromString(body, &v); err != nil {
			return
		}
		ip = v["query"].(string)
	default:
		_, body, errs := gorequest.New().Get("https://checkip.amazonaws.com").End()
		if len(errs) > 0 {
			return nil, errs[0]
		}
		ip = body
	}
	return net.ParseIP(ip), nil
}

type Proxy struct {
	Type string
	Host string
	Port int
}

func (p Proxy) ApplyGorequest(req *gorequest.SuperAgent) error {
	if p.Host != "" {
		pu := fmt.Sprintf("%s://%s:%d", p.Type, p.Host, p.Port)
		switch p.Type {
		case "socks5":
			u, e := url.Parse(pu)
			if e != nil {
				return e
			}
			socks5Dialer, e1 := proxy.FromURL(u, proxy.Direct)
			if e1 != nil {
				return e
			}
			req.Transport = &http.Transport{Dial: socks5Dialer.Dial}
		case "http", "https":
			req.Proxy(pu)
		}
	}

	return nil
}
