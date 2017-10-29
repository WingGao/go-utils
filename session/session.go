package session

import (
	"encoding/gob"
	"bytes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/kataras/iris/context"
	. "utils"
	//"fmt"
	"errors"
	cdb "core/db"
	pxdb "px/db"
	"time"
	"net/http/httptest"
	"net/http"
)

type XSession struct {
	ctx      context.Context
	key      string
	isClear  bool
	Sid      string
	Uid      uint32
	Group    uint32
	Username string
	LastTime time.Time
}

var (
	_session   *sessions.Sessions
	_errNotSet = errors.New("utils.session not set")
	_rdb       *redis.Database
)

func BuildIrisSession(conf MConfig) {
	rconf := conf.Redis
	_rdb = redis.New(service.Config{
		Network:     service.DefaultRedisNetwork,
		Addr:        rconf.Addr,
		Password:    rconf.Password,
		Database:    rconf.Database,
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      ""})

	_rdb.Async(true)
	iris.RegisterOnInterrupt(func() {
		_rdb.Close()
	})

	exp := time.Duration(conf.CookieExpires)
	if exp > 0 {
		exp = exp * time.Second
	}
	mySessions := sessions.New(sessions.Config{
		Cookie:  "smsid",
		Expires: exp,
	})
	mySessions.UseDatabase(_rdb)

	SetIrisSessionFactory(mySessions)
	return
}

func SetIrisSessionFactory(s *sessions.Sessions) {
	if _session != nil {
		return
	}
	_session = s
}

func checkSession() bool {
	return _session != nil
}

func GetIrisSession(ctx context.Context) *sessions.Session {
	sess := _session.Start(ctx)
	return sess
}
func getSessionCtx(key string) context.Context {
	ctx := context.NewContext(nil)
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.AddCookie(&http.Cookie{Name: "smsid", Value: key})
	w := httptest.NewRecorder()
	ctx.BeginRequest(w, req)
	return ctx
}

func GetIrisSessionByKey(key string) *sessions.Session {
	ctx := getSessionCtx(key)
	defer func() {
		ctx.EndRequest()
	}()
	return GetIrisSession(ctx)
}
func NewSessionByKey(key string) (*XSession, error) {
	ctx := getSessionCtx(key)
	defer func() {
		ctx.EndRequest()
	}()
	return NewSessionFromIris(ctx, XSESSION_KEY)
}
func NewSessionFromIris(ctx context.Context, key string) (*XSession, error) {
	if !checkSession() {
		return nil, _errNotSet
	}
	sess := _session.Start(ctx)
	val := sess.Get(key)
	if val == nil {
		return &XSession{ctx: ctx, key: key}, nil
	}
	xsess, err := NewSessionFromGob(val.([]byte))
	xsess.RefreshAuto()
	xsess.ctx = ctx
	xsess.key = key
	return xsess, err
}

func NewSessionFromGob(bs []byte) (*XSession, error) {
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)
	sess := &XSession{}
	err := dec.Decode(&sess)
	return sess, err
}

func (x *XSession) ToGob() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(x)
	return buf.Bytes(), err
}

func (x *XSession) Clear() () {
	x.isClear = true
	x.Uid = 0
	x.Group = 0
}

func (x *XSession) SaveIris(ctx context.Context, key string) error {
	if !checkSession() {
		return _errNotSet
	}
	//被清空的不需要保存
	if x.isClear {
		_session.Destroy(ctx)
		return nil
	}

	g, err := x.ToGob()
	sess := _session.Start(ctx)
	sess.Set(key, g)
	return err
}

// 直接保存
func (x *XSession) SaveIrisD() error {
	if x.ctx == nil || x.key == "" {
		return errors.New("XSession ctx or key is not set")
	}

	g, err := x.ToGob()
	sess := _session.Start(x.ctx)
	sess.Set(x.key, g)
	return err
}

func (x *XSession) IsValued() bool {
	return x.Uid > 0
}
func IsAdmin(group uint32) bool {
	return group == cdb.GROUP_ADMIN
}
func (x *XSession) IsAdmin() bool {
	return IsAdmin(x.Group)
}

func (x *XSession) IsTeacher() bool {
	return x.Group == cdb.GROUP_TEACHER
}
func (x *XSession) Teacher() (s *pxdb.Teacher) {
	if x.IsTeacher() {
		s = pxdb.NewTeacher()
		s.LoadAndSetId(x.Uid)
	}
	return
}
func (x *XSession) IsStudent() bool {
	return x.Group == cdb.GROUP_STUDENT
}
func (x *XSession) Student() (s *pxdb.Student) {
	if x.IsStudent() {
		s = pxdb.NewStudent()
		s.LoadAndSetId(x.Uid)
	}
	return
}
func (x *XSession) Account() (s *cdb.Account) {
	if x.Uid > 0 {
		s = cdb.NewAccount()
		err := s.LoadAndSetId(x.Uid)
		if err != nil {
			return nil
		}
	}
	return
}
func (x *XSession) Refresh() {
	acc := x.Account()
	if acc == nil {
		//清空
		x.Clear()
	} else {
		x.Group = acc.Group
		x.LastTime = time.Now()
	}
}

//自动session检查
func (x *XSession) RefreshAuto() {
	//5分钟检查一次session
	//用户可能被删除、更新的情况
	if x.Uid > 0 && time.Now().After(x.LastTime.Add(5*time.Minute)) {
		x.Refresh()
	}
}
