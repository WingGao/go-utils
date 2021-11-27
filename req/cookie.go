package req

import (
	"github.com/WingGao/go-utils/ucore"
	"net/http"
	"net/url"
	"strings"
)

// ParseHeaderCookies 解析Http请求header中的cookie
func ParseHeaderCookies(headerCookies string) []*http.Cookie {
	header := http.Header{}
	header.Add("Cookie", headerCookies)
	r := http.Request{Header: header}
	cookies := r.Cookies()
	return cookies
}

func AddCookieToJar(jar http.CookieJar, headerCookies string, url1 string) {
	cookies := ParseHeaderCookies(headerCookies)
	u, _ := url.Parse(url1)
	jar.SetCookies(u, cookies)
}

func ParseHeaderCookiesToMap(headerCookies string) map[string]string {
	cmap := make(map[string]string)
	for _, pair := range strings.Split(headerCookies, ";") {
		pp := strings.Split(pair, "=")
		cmap[strings.TrimSpace(pp[0])] = strings.TrimSpace(pp[1])
	}
	return cmap
}

func JoinCookieMapToString(hmap map[string]string) string {
	sb := ucore.StringBuilder{}
	cnt := 0
	for k, v := range hmap {
		if cnt > 0 {
			sb.Write("; ")
		}
		sb.Write(k, "=", v)
		cnt++
	}
	return sb.String()
}
