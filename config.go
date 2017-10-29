package utils

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
	tconfig "github.com/RichardKnop/machinery/v1/config"
)

//main config
type MConfig struct {
	Debug           bool
	Addr            string //服务地址
	Host            string
	Mysql           string
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

	LibreOfficePath    string `yaml:"libreoffice_path"`
	TestAdminSession   string
	TestStudentSession string
	TestTeacherSession string
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
}

func LoadConfig(confPath string) (MConfig, error) {
	conf := MConfig{
		Debug: false,
	}
	if confPath == "" {

	} else {
		confd, _ := ioutil.ReadFile(confPath)
		err := yaml.Unmarshal(confd, &conf)
		if err != nil {
			return conf, err
		} else {
			fmt.Println("MConfig loaded")
		}
	}
	return conf, nil
}
