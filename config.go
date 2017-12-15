package utils

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	tconfig "github.com/RichardKnop/machinery/v1/config"
	"path/filepath"
	"github.com/go-errors/errors"
)

var (
	DefaultConfig MConfig
)
//main config
type MConfig struct {
	configPath      string                 //配置文件路径
	allMap          map[string]interface{} //保存配置,默认map[interface {}]interface {}
	AppPath         string                 //运行路径，一般不设置，测试使用
	Debug           bool
	Addr            string //服务地址
	Host            string
	Mysql           string
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
}

func (m MConfig) GetConfigPath() string {
	return m.configPath
}

func (m MConfig) Get(key string) interface{} {
	if v, ok := m.allMap[key]; ok {
		return v
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

type WxConfig struct {
	AppId     string
	MchId     string
	ApiKey    string
	NotifyUrl string
	Token     string //缓存用
}

type RedisConf struct {
	Addr        string
	Password    string
	Database    string
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
			m := make(map[string]interface{})
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
