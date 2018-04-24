package session

import (
	"testing"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
	"mtest"
	"os"
	uredis "github.com/WingGao/go-utils/redis"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mtest.PreEnv()
	testConf := mtest.GetConfig()
	uredis.LoadClient(testConf.Redis)
	BuildIrisSession(testConf)
	os.Exit(m.Run())
}
func TestXSession_ToGob(t *testing.T) {
	sess := &XSession{Uid: 123}
	g, _ := sess.ToJSON()
	t.Logf("session: %s", g)
	s2, _ := NewSessionFromJSON(g)
	t.Logf("%#v", s2)
}

func TestSessionFromIris(t *testing.T) {
	val := XSession{
		Username: "testname",
	}
	app := iris.New()
	app.Post("/set", func(ctx context.Context) {
		sess, _ := NewSessionFromIris(ctx, "sess")
		sess.Uid = val.Uid
		if err := sess.SaveIrisD(); err != nil {
			t.Error(err)
			ctx.StatusCode(403)
		}
	})
	app.Get("/get", func(ctx context.Context) {
		sess, err := NewSessionFromIris(ctx, "sess")
		if err != nil {
			t.Error(err)
			ctx.StatusCode(403)
			return
		}
		if sess.Uid == val.Uid {
			ctx.JSON(val)
		} else {
			ctx.JSON(sess)
		}
	})

	e := httptest.New(t, app, httptest.URL("http://example.com"))
	e.POST("/set").Expect().Status(iris.StatusOK).Cookies().NotEmpty()
	e.GET("/get").Expect().Status(iris.StatusOK).JSON().Object().Equal(val)
}

func TestGetIrisSessionByKey(t *testing.T) {
	k := "f02f5ad1-7994-46fd-aa9a-301f74ca869c"
	sess := GetIrisSessionByKey(k)
	smap := sess.GetAll()
	assert.NotEmpty(t, smap)
	xsess, err := NewSessionByKey(k)
	assert.NoError(t, err)
	t.Log(xsess)
}

func TestClearUserAllSessions(t *testing.T) {
	err := ClearUserAllSessions(169)
	assert.NoError(t, err)
}
