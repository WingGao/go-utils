package utils

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	tconfig "github.com/RichardKnop/machinery/v1/config"
	"path/filepath"
	"github.com/go-errors/errors"
	"strings"
)

var (
	DefaultConfig MConfig
)
//main config
type MConfig struct {
	configPath string                      //配置文件路径
	allMap     map[interface{}]interface{} //保存配置,默认map[interface {}]interface {}
	AppPath    string                      //运行路径，一般不设置，测试使用
	Debug      bool
	AdminMail  string
	Project    string //工程名字
	Host       string
	Port       string //服务地址
	PublicHost string //对外的部署地址 host:port
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		Option   string
	}
	MysqlDebug      bool
	Mongodb         string
	DefaultPassword string `yaml:"default_password"` //默认密码
	MediaPath       string
	WebApps         string
	MaxMediaSize    string `yaml:"max_media_size"`
	CookieExpires   int64  `yaml:"cookie_expires"`
	SiteCreator struct {
		Mysql string
		WpDir string
	}
	WxConfig WxConfig `yaml:"wechat"`
	Redis    RedisConf
	Task     tconfig.Config

	Cms struct {
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
	}
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

func (m MConfig) GetMySQLString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.Mysql.User, m.Mysql.Password,
		m.Mysql.Host, m.Mysql.Port, m.Mysql.DBName, m.Mysql.Option)
}

func (m MConfig) GetMachineryConfig() *tconfig.Config {
	scnf := &m.Task
	if scnf.Broker == "" { //默认使用redis
		scnf.Broker = fmt.Sprintf("redis://%s/%d", m.Redis.Addr, m.Redis.Database)
	}
	return scnf
}

type WxConfig struct {
	AppId     string
	MchId     string
	MchApiKey    string
	NotifyUrl string
	Token     string //缓存用
}

type RedisConf struct {
	Addr        string //host:port 127.0.0.1:6379
	Password    string
	Database    int
	UniqueIdKey string
	Prefix      string
}

func NewConfigFromFile(confPath string) (conf MConfig, err error) {
	if confPath == "" {
		err = errors.New("need path")
	} else {
		confPath, _ = filepath.Abs(confPath)
		confd, _ := ioutil.ReadFile(confPath)
		err := yaml.Unmarshal(confd, &conf)
		if err != nil {
			return conf, err
		} else {
			m := make(map[interface{}]interface{})
			yaml.Unmarshal(confd, &m)
			conf.allMap = m
			conf.configPath = confPath
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
		fmt.Println("MConfig loaded")
	}
	fmt.Println("AppPath:", conf.AppPath)
	fmt.Println("MediaPath:", conf.MediaPath)
	fmt.Println("WebApps:", conf.WebApps)
	DefaultConfig = conf
	return conf, err
}
