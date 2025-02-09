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
	TencentCDN TencentCDN `yaml:"TencentCDN"`
}

type TencentCDN struct {
	CDNDomain string             `yaml:"CDNDomain"`
	Region    string             `yaml:"Region"`
	Bucket    string             `yaml:"Bucket"`
	SecretID  string             `yaml:"SecretID"`
	SecretKey string             `yaml:"SecretKey"`
	Tenants   []TencentCDNTenant `yaml:"Tenants"`
}

type TencentCDNTenant struct {
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
