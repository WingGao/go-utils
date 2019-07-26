package session

import (
	"github.com/WingGao/go-utils/session/redis"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/sessions"
	sredis "github.com/kataras/iris/sessions/sessiondb/redis"

	//"fmt"
	"errors"
	"fmt"
	uredis "github.com/WingGao/go-utils/redis"
	"github.com/chanxuehong/wechat/oauth2"
	"github.com/jinzhu/copier"
	"github.com/json-iterator/go"
	"net/http"
	"net/http/httptest"
	//pxdb "px/db"
	"time"
)
const (
	XSESSION_KEY = "xsession"
)
type XSession struct {
	ctx      context.Context
	Iris     *sessions.Session `json:"-"`
	key      string //保存在iris中的键
	isClear  bool
	parent   interface{}
	Sid      string
	Uid      uint32
	Group    uint32
	Username string
	LastTime time.Time
	WxOpenId string
	WxToken  *oauth2.Token
	Items    map[string]interface{}
}
type IXSession interface {
	IsClear() bool
	Clear() ()
	IsValued() bool
	Refresh()
	RefreshAuto()
	New() interface{}
}

var (
	_session          *sessions.Sessions
	_errNotSet        = errors.New("utils.session not set")
	_rdb              *redis.Database
	_sessionKeyPrefix = "core:sid:"
)

func BuildIrisSession(rconf uredis.RedisConf,cookieExpire int64) {
	_rdb = redis.New(uredis.MainClient, sredis.Config{
		Network:     sredis.DefaultRedisNetwork,
		Addr:        rconf.Addr,
		Password:    rconf.Password,
		Database:    fmt.Sprintf("%d", rconf.Database),
		MaxActive:   0,
		Prefix:      _sessionKeyPrefix})

	iris.RegisterOnInterrupt(func() {
		_rdb.Close()
	})

	exp := time.Duration(cookieExpire)
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
	xsess, err := NewSessionFromJSON(val.(string))
	xsess.RefreshAuto()
	xsess.ctx = ctx
	xsess.key = key
	xsess.Iris = sess
	xsess.Sid = sess.ID()
	return xsess, err
}

func NewSessionFromJSON(bs string) (*XSession, error) {
	//buf := bytes.NewBuffer(bs)
	//dec := gob.NewDecoder(buf)
	sess := &XSession{}
	//err := dec.Decode(&sess)
	err := jsoniter.UnmarshalFromString(bs, sess)
	return sess, err
}

func (x *XSession) ToJSON() (string, error) {
	//var buf bytes.Buffer
	//enc := gob.NewEncoder(&buf)
	////make a copy
	//dx := &XSession{}
	//copier.Copy(dx, x)
	//dx.ctx = nil
	//dx.Iris = nil
	//err := enc.Encode(dx)
	//return buf.Bytes(), err
	return jsoniter.MarshalToString(x)
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

	g, err := x.ToJSON()
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

func (x *XSession) IsValid() bool {
	return x.Uid > 0
}
func (x *XSession) Refresh() {
	//TODO 测试
}

//自动session检查
func (x *XSession) RefreshAuto() {
	//5分钟检查一次session
	//用户可能被删除、更新的情况
	if x.Uid > 0 && time.Now().After(x.LastTime.Add(5*time.Minute)) {
		x.Refresh()
	}
}

func (x *XSession) Set(key string, val interface{}) {
	if x.Items == nil {
		x.Items = make(map[string]interface{}, 50)
	}
	x.Items[key] = val
}

func (x *XSession) Get(key string) (val interface{}, ok bool) {
	if x.Items == nil {
		return nil, false
	}
	val, ok = x.Items[key]
	return
}
func (x *XSession) GetString(key string) (val string, ok bool) {
	if v1, o1 := x.Get(key); !o1 {
		return "", false
	} else {
		return v1.(string), true
	}
}

//删除所有的用户登录session
func ClearUserAllSessions(uid uint32) (err error) {
	userKey := fmt.Sprintf("core:user:sids:%d", uid)
	sids, err2 := uredis.MainClient.SMembers(userKey).Result()
	if err2 != nil {
		return err2
	}
	for _, v := range sids {
		_session.DestroyByID(v)
		//手动调用，可能不在内存里
		_rdb.Release(v)
	}
	//删除自己
	uredis.MainClient.Del(userKey)
	return
}

//记录到登录列表,目前我们把登录列表放在redis中
func AddUserLoginSession(uid uint32, sid string) error {
	userKey := fmt.Sprintf("core:user:sids:%d", uid)
	_, err := uredis.MainClient.SAdd(userKey, sid).Result()
	return err
}
