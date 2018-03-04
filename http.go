package utils

import (
	"github.com/parnurzeal/gorequest"
	"regexp"
)

// 创建一个gorequest, 默认是电脑端的UA
func NewGorequest() (r *gorequest.SuperAgent) {
	r = gorequest.New()
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	return
}

func GetPublicIP() (string, error) {
	_, body, errs := gorequest.New().Get("http://www.net.cn/static/customercare/yourip.asp").End()
	if len(errs) > 0 {
		return "", errs[0]
	}
	re, _ := regexp.Compile(`<h2>([\d\.]+)</h2>`)
	ip := re.FindSubmatch([]byte(body))
	return string(ip[1]), nil
}
