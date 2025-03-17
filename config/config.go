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
	ListenPort         string       `yaml:"ListenPort"`
	MaxRequestBodySize int          `yaml:"MaxRequestBodySize"`
	Matomo             MatomoClient `yaml:"Matomo"`
	Services           Services     `yaml:"Services"`
}

type Services struct {
	RailgunCDN RailgunCDN `yaml:"RailgunCDN"`
}

type RailgunCDN struct {
	COS     TencentCOS                                 `yaml:"COS"`
	CDN     TencentCDN                                 `yaml:"CDN"`
	Private PrivateCDN                                 `yaml:"Private"`
	Tenants map[RailgunCDNTenantAppID]RailgunCDNTenant `yaml:"Tenants"`
}

type RailgunCDNTenantAppID = string

type RailgunCDNTenant struct {
	AppKey   string `yaml:"AppKey"`
	RootPath string `yaml:"RootPath"`
	SiteID   string `yaml:"SiteID"`
}

type TencentCOS struct {
	Region    string `yaml:"Region"`
	Bucket    string `yaml:"Bucket"`
	SecretID  string `yaml:"SecretID"`
	SecretKey string `yaml:"SecretKey"`
}

type TencentCDN struct {
	Endpoint        string `yaml:"Endpoint"`
	PKey            string `yaml:"PKey"`
	TimestampOffset int64  `yaml:"TimestampOffset"`
}

type PrivateCDN struct {
	Endpoint string `yaml:"Endpoint"`
}

type MatomoClient struct {
	Endpoint        string `yaml:"Endpoint"`
	AuthToken       string `yaml:"AuthToken"`
	NumWorkers      int    `yaml:"NumWorkers"`
	BatchSize       int    `yaml:"BatchSize"`
	EventBufferSize int    `yaml:"EventBufferSize"`
}

func Init() {
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		hlog.Fatalf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", &Conf); err != nil {
		hlog.Fatalf("error unmarshalling config: %v", err)
	}
}
