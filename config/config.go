package config

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	k    = koanf.New(".")
	path = "config.yaml"
	Conf Config
)

type Config struct {
	ListenPort string   `yaml:"ListenPort"`
	Services   Services `yaml:"Services"`
}

type Services struct {
	RailgunCDN RailgunCDN `yaml:"RailgunCDN"`
}

type RailgunCDN struct {
	COS     TencentCOS         `yaml:"COS"`
	CDN     TencentCDN         `yaml:"CDN"`
	Tenants []RailgunCDNTenant `yaml:"Tenants"`
}

type TencentCOS struct {
	Region    string `yaml:"Region"`
	Bucket    string `yaml:"Bucket"`
	SecretID  string `yaml:"SecretID"`
	SecretKey string `yaml:"SecretKey"`
}

type TencentCDN struct {
	Domain string `yaml:"Domain"`
	PKey   string `yaml:"PKey"`
}

type RailgunCDNTenant struct {
	RootPath string `yaml:"RootPath"`
	AppID    string `yaml:"AppID"`
	AppKey   string `yaml:"AppKey"`
}

func Init() {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		hlog.Fatalf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", &Conf); err != nil {
		hlog.Fatalf("error unmarshalling config: %v", err)
	}
}
