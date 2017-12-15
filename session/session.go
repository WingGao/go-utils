package session

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
	"github.com/kataras/iris/context"
	. "github.com/WingGao/go-utils"
	//"fmt"
	"errors"
	cdb "core/db"
	uredis "github.com/WingGao/go-utils/redis"
	//pxdb "px/db"
	"time"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	"github.com/json-iterator/go"
	"fmt"
	"github.com/jinzhu/copier"
)

type XSession struct {
	ctx      context.Context
	Iris     *sessions.Session `json:"-"`
	key      string //保存在iris中的键
	isClear  bool
	Sid      string
	Uid      uint32
	Group    uint32
	Username string
	LastTime time.Time
	WxOpenId string
}
type IXSession interface {
	IsClear() bool
	Clear() ()
	IsValued() bool
	Refresh()
	RefreshAuto()
}

var (
	_session          *sessions.Sessions
	_errNotSet        = errors.New("utils.session not set")
	_rdb              *redis.Database
	_sessionKeyPrefix = "core_sid_"
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
		Prefix:      _sessionKeyPrefix})

	_rdb.Async(true)
	iris.RegisterOnInterrupt(func() {
		_rdb.Close()
	})

	exp := time.Duration(conf.CookieExpires)
	if exp > 0 {
		exp = exp * time.Second
	}
	mySessions := sessions.New(sessions.Config{
		Cookie:       "smsid",
		Expires:      exp,
		AllowReclaim: true,
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
func GetSessionCtx(key string) context.Context {
	ctx := context.NewContext(nil)
	req := httptest.NewRequest("GET", "http://localhost", nil)
	req.AddCookie(&http.Cookie{Name: "smsid", Value: key})
	w := httptest.NewRecorder()
	ctx.BeginRequest(w, req)
	return ctx
}

func GetIrisSessionByKey(key string) *sessions.Session {
	ctx := GetSessionCtx(key)
	defer func() {
		ctx.EndRequest()
	}()
	return GetIrisSession(ctx)
}

func NewSessionByKey(key string) (*XSession, error) {
	ctx := GetSessionCtx(key)
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
		return &XSession{ctx: ctx, key: key, Iris: sess}, nil
	}
	xsess, err := NewSessionFromGob(val.([]byte))
	xsess.RefreshAuto()
	xsess.ctx = ctx
	xsess.key = key
	xsess.Iris = sess
	xsess.Sid = sess.ID()
	return xsess, err
}

func NewSessionFromGob(bs []byte) (*XSession, error) {
	//buf := bytes.NewBuffer(bs)
	//dec := gob.NewDecoder(buf)
	sess := &XSession{}
	//err := dec.Decode(&sess)
	err := jsoniter.Unmarshal(bs, sess)
	return sess, err
}

func (x *XSession) ToGob() ([]byte, error) {
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	////make a copy
	//dx := &XSession{}
	//copier.Copy(dx, x)
	//dx.ctx = nil
	//dx.Iris = nil
	//err := enc.Encode(dx)
	//return buf.Bytes(), err
	return json.Marshal(x)
}

func (x *XSession) UpdateExpiration(expires time.Duration) {
	_session.UpdateExpiration(x.ctx, expires)
}

func (x *XSession) IsClear() bool {
	return x.isClear
}

func (x *XSession) Clear() {
	newSess := XSession{ctx: x.ctx, key: x.key}
	copier.Copy(x, newSess)
	x.isClear = true
}

func (x *XSession) SaveIris(ctx context.Context, key string) error {
	if !checkSession() {
		return _errNotSet
	}
	//被清空的不需要保存
	if x.isClear {
		_session.Destroy(ctx)
		x.isClear = false
	}

	g, err := x.ToGob()
	sess := x.Iris
	if sess == nil {
		sess = _session.Start(ctx)
		x.Iris = sess
	}
	sess.Set(key, g)
	return err
}

// 直接保存
func (x *XSession) SaveIrisD() error {
	if x.ctx == nil || x.key == "" {
		return errors.New("XSession ctx or key is not set")
	}

	return x.SaveIris(x.ctx, x.key)
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

//删除所有的用户登录session
func ClearUserAllSessions(uid uint32) (err error) {
	userKey := fmt.Sprintf("core_user_%d_sids", uid)
	sids := []string{}
	err = uredis.MainClient.Csmembers(userKey, &sids)
	if err != nil {
		return
	}
	for i, v := range sids {
		_session.DestroyByID(v)
		if i == 0 { //初始化
			_rdb.Load("")
		}
		//手动调用，可能不在内存里
		_rdb.Sync(sessions.SyncPayload{
			SessionID: v,
			Action:    sessions.ActionDestroy,
		})
	}
	//删除自己
	uredis.MainClient.Del(userKey)
	return
}
