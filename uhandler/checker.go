package uhandler

import (
	"github.com/kataras/iris/context"
	cdb "core/db"
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

func RequireAdmin(ictx context.Context) {
	sess, _ := NewSessionFromIris(ictx, XSESSION_KEY)
	if !sess.IsAdmin() {
		panic(ERR_REQUIRE_ADMIN)
	}
}

//PX
func RequireTeacher(sess *XSession) {
	err := cdb.CheckGroup(cdb.GROUP_TEACHER, sess.Group)
	if err != nil {
		panic("need teacher")
	}
}

