package utils

import (
	"encoding/hex"
	"fmt"
	tconfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/WingGao/go-utils/redis"
	"github.com/go-errors/errors"
	"github.com/micro/go-micro/registry"
	"github.com/ungerik/go-dry"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	DefaultConfig MConfig
)

type DbConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	Option      string
	CreateDB    bool                  //创建对应数据库
	NotCheck    bool                  //是否检查
	StartCmd    string `yaml:"start"` //启动命令
	AutoMigrate bool                  //自动建表
}

type ThirdPartConfig struct {
	Name   string
	AppID  string
	Key    string
	Secret string
}

//main config
type MConfig struct {
	configPath      string                      //配置文件路径
	allMap          map[interface{}]interface{} //保存配置,默认map[interface {}]interface {}
	AppPath         string                      //运行路径，一般不设置，测试使用
	Debug           bool
	UID             int //linux uid
	GID             int //linux gid
	AdminMail       string
	Project         string //工程名字
	Host            string
	Port            string //服务地址
	Https           bool
	PublicHost      string //对外的部署地址 host:port
	Grpc            GrpcConfig
	MasterKey       string
	Mysql           DbConfig
	MysqlDebug      bool
	AutoBackup      bool //自动备份
	Postgresql      DbConfig
	ElasticSearch   DbConfig `yaml:"elastic"`
	Mongodb         DbConfig
	DefaultPassword string `yaml:"default_password"` //默认密码
	CasbinModelPath string `yaml:"casbin_model"`
	CasbinPolicy    string `yaml:"casbin_policy"`
	MediaPath       string
	Media           struct {
		Type   string
		Bucket string
		Region string
		Host   string
	}
	WebApps        string
	SsrMap         map[string]string `yaml:"ssr"`
	MaxMediaSize   string            `yaml:"max_media_size"`
	CookieExpires  int64             `yaml:"cookie_expires"`
	CaptchaDisable bool              `yaml:"captcha_disable"`
	SiteCreator    struct {
		Mysql string
		WpDir string
	}
	WxConfig WxConfig `yaml:"wechat"`
	Redis    redis.RedisConf
	Task     tconfig.Config
	Seo      bool
	Cms      struct {
	}

	YunCourseTemplate string

	LibreOfficePath    string `yaml:"libreoffice_path"`
	TestAdminSession   string
	TestStudentSession string
	TestTeacherSession string

	Wing struct {
		FaceKey string
		FaceSec string
	}

	Exam struct {
		Mysql string
	}

	Smtp struct {
		Server   string
		Port     int
		SSL      bool
		User     string
		Password string
		From     string
	}
	Qiniu  ThirdPartConfig //七牛
	Aliyun struct {
		ThirdPartConfig     `yaml:",inline"`
		RegionId            string
		SmsSign             string `yaml:"sms_sign"`
		SmsTemplateRegister string `yaml:"sms_template_register"`
	}
	Vips    string //vips路径， '-'表示不需要
	GraphQL GraphQLConf
	Kafka   KafkaConfig
	Log     LogConfig
}

func (m MConfig) GetConfigPath() string {
	return m.configPath
}

func (m MConfig) Get(key string) interface{} {
	last := m.allMap
	keys := strings.Split(key, ".")
	for i, k := range keys {
		if v, ok := last[k]; ok {
			if i+1 == len(keys) { //最后一个
				return v
			} else { //下一级
				if last, ok = v.(map[interface{}]interface{}); !ok { //转换失败，没有子元素
					return nil
				}
			}
		} else {
			return nil
		}
	}

	return nil
}

func (m MConfig) GetString(key, def string) string {
	v := m.Get(key)
	if v == nil || v.(string) == "" {
		return def
	}
	return v.(string)
}

//获得相对于配置文件的绝对路径
func (m MConfig) AbsPath(apath string) string {
	return getFullPath(filepath.Dir(m.configPath), apath)
}

func (m MConfig) GetMySQLString(dbname string) string {
	if dbname == "" {
		dbname = m.Mysql.DBName
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.Mysql.User, m.Mysql.Password,
		m.Mysql.Host, m.Mysql.Port, dbname, m.Mysql.Option)
}

func (m MConfig) GetPostgresqlString(dbname string) string {
	if dbname == "" {
		dbname = m.Postgresql.DBName
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", m.Postgresql.User, m.Postgresql.Password,
		m.Postgresql.Host, m.Postgresql.Port, dbname, m.Postgresql.Option)
}

func (m MConfig) GetMachineryConfig() *tconfig.Config {
	scnf := &m.Task
	if scnf.Broker == "" { //默认使用redis
		scnf.Broker = fmt.Sprintf("redis://%s", redis.MainClient.GetConfig().Addr)
	}
	if scnf.ResultBackend == "" { //默认使用redis
		scnf.ResultBackend = scnf.Broker
	}
	return scnf
}

// 拼接处完整的对外地址，http协议
func (m MConfig) BuildUrl(relpath string) string {
	if !strings.HasPrefix(relpath, "/") {
		relpath = "/" + relpath
	}
	prototal := "http"
	if m.Https {
		prototal = "https"
	}
	return fmt.Sprintf("%s://%s%s", prototal, m.PublicHost, relpath)
}

//创建目录，并制定默认uid，gid
func (m MConfig) MkdirAll(fp string, mod os.FileMode) (err error) {
	err = os.MkdirAll(fp, mod)
	if err != nil {
		return
	}
	if m.UID > 0 || m.GID > 0 {
		err = os.Chown(fp, m.UID, m.GID)
	}
	return
}

//创建目录，并制定默认uid，gid
func (m MConfig) MkdirAllDef(fp string) (err error) {
	return m.MkdirAll(fp, DEFAULT_FILEMODE)
}

func (m MConfig) EncryptToHex(plain string) string {
	key := []byte(m.MasterKey)
	outb := dry.EncryptAES(key, []byte(plain))
	outs := hex.EncodeToString(outb)
	return outs
}

func (m MConfig) DecryptFromHex(sec string) string {
	key := []byte(m.MasterKey)
	secb, _ := hex.DecodeString(sec)
	outb := dry.DecryptAES(key, secb)
	return string(outb)
}

type WxConfig struct {
	AppId             string
	AppKey            string
	MchId             string
	MchApiKey         string
	NotifyUrl         string
	EnableTokenServer bool `yaml:"enable_tokenserver"`
	EnableJsTicket    bool `yaml:"enable_jsticket"`
	Token             string //缓存用
	Miniapps          []*ThirdPartConfig `yaml:"mini"`
	miniappsMap       map[string]*ThirdPartConfig //通过Miniapps转换的
	Corp              WxCorpConf
}

func (m *WxConfig) Update() error {
	m.miniappsMap = make(map[string]*ThirdPartConfig, len(m.Miniapps))
	for _, c := range m.Miniapps {
		if _, exist := m.miniappsMap[c.Name]; exist {
			return errors.Errorf("小程序%s[%s]已存在", c.Name, c.AppID)
		}
		m.miniappsMap[c.Name] = c
	}
	return nil
}

func (m *WxConfig) GetMiniApp(name string) (*ThirdPartConfig, bool) {
	if name == "" {
		name = "default"
	}
	app, ok := m.miniappsMap[name]
	return app, ok
}

type WxCorpConf struct {
	CorpId        string
	CorpSecret    string
	OauthRedirect string `yaml:"oauth_redirect"`
}

type GraphQLConf struct {
	Path string
}

type GrpcConfig struct {
	Port        int
	ServicesMap map[string][]registry.Node `yaml:"services"`
}

type KafkaConfig struct {
	AppId     string
	Addresses string
}
type LogConfig struct {
	NoStd   bool
	Kafka   bool
	Request bool
}

func NewConfigFromFile(confPath string) (conf MConfig, err error) {
	if confPath == "" {
		err = errors.New("need path")
	} else {
		confPath, _ = filepath.Abs(confPath)
		confd, err1 := ioutil.ReadFile(confPath)
		if err1 != nil {
			err = err1
			return
		}
		err = yaml.Unmarshal(confd, &conf)
		if err != nil {
			return conf, err
		} else {
			m := make(map[interface{}]interface{})
			yaml.Unmarshal(confd, &m)
			conf.allMap = m
			conf.configPath = confPath
			// 默认已配置文件所在目录为准
			if conf.AppPath == "" {
				conf.AppPath = filepath.Dir(confPath)
			}
			if conf.MediaPath == "" {
				conf.MediaPath = filepath.Join(conf.AppPath, "uploads")
			} else {
				conf.MediaPath = getFullPath(conf.AppPath, conf.MediaPath)
			}
			if conf.WebApps == "" {
				conf.WebApps = filepath.Join(conf.AppPath, "webapps")
			} else {
				conf.WebApps = getFullPath(conf.AppPath, conf.WebApps)
			}
			if conf.GraphQL.Path != "" {
				conf.GraphQL.Path = getFullPath(conf.AppPath, conf.GraphQL.Path)
			}
			if conf.CasbinModelPath != "" {
				conf.CasbinModelPath = getFullPath(conf.AppPath, conf.CasbinModelPath)
			}
			conf.Mysql.Host = formatEnv(conf.Mysql.Host)
			err = conf.WxConfig.Update()
		}
	}
	return
}

func getFullPath(apppath, p string) string {
	if p == "" {
		return apppath
	}
	if filepath.IsAbs(p) {
		return p
	} else {
		return filepath.Join(apppath, p)
	}
}
func LoadConfig(confPath string) (MConfig, error) {
	conf, err := NewConfigFromFile(confPath)
	if err == nil {
		fmt.Printf("MConfig loaded %s\n", confPath)
	}
	fmt.Println("AppPath:", conf.AppPath)
	fmt.Println("MediaPath:", conf.MediaPath)
	fmt.Println("WebApps:", conf.WebApps)
	DefaultConfig = conf
	return conf, err
}

func formatEnv(v string) string {
	if strings.HasPrefix(v, "$") {
		return dry.GetenvDefault(string([]rune(v)[1:]), "")
	}
	return v
}
