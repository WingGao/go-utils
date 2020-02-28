package whandler

import (
	ucore "github.com/WingGao/go-utils/ucore"
	. "github.com/WingGao/go-utils/session"
)

func RequireLoginX(sess *XSession) {
	if sess.Uid <= 0 {
		panic(ucore.NewErrNeedLogin())
	}
}
