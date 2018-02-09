package utils

import (
	"github.com/parnurzeal/gorequest"
)

// 创建一个gorequest, 默认是电脑端的UA
func NewGorequest() (r *gorequest.SuperAgent) {
	r = gorequest.New()
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	return
}
