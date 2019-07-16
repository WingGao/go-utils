package uhandler

import (
	"github.com/kataras/iris/context"
	. "github.com/WingGao/go-utils"
	. "github.com/WingGao/go-utils/session"
)

func RequireLogin(ictx context.Context) {
	sess, _ := NewSessionFromIris(ictx, XSESSION_KEY)
	RequireLoginX(sess)
}

func RequireLoginX(sess *XSession) {
	if sess.Uid <= 0 {
		panic(ERR_REQUIRE_LOGIN)
	}
}
