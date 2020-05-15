package ucore

import (
	"github.com/parnurzeal/gorequest"
	"net"
	"regexp"
)

func GetRealIP() net.IP {
	_, body, errs := gorequest.New().Get("https://2020.ip138.com/").End()
	if len(errs) > 0 {
		panic(errs[0])
	}
	r, _ := regexp.Compile(`\[([\d+\.]+)\]`)
	ips := r.FindStringSubmatch(body)
	return net.ParseIP(ips[1])
}
