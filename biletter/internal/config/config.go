package config

import (
	"biletter/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug    *bool         `yaml:"is_debug" env-required:"true"`
	Listen     listen        `yaml:"listen"`
	Storage    StorageConfig `yaml:"storage"`
	Redis      RedisConfig   `yaml:"redis"`
	Grpc       Grpc          `yaml:"grpc"`
	GrpcClient GrpcClient    `yaml:"grpcClient"`
	Cors       Cors          `yaml:"cors"`
	MainURI    string        `yaml:"mainURI"`
	Timezone   string        `yaml:"timezone"`
}

type listen struct {
	Type   string `yaml:"type" env-default:"port"`
	BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
	Port   string `yaml:"port" env-default:"3000"`
}

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Grpc struct {
	BindIP string `yaml:"bind_ip"`
	Port   string `yaml:"port"`
	Ssl    bool   `yaml:"ssl"`
}

type FileStorage struct {
	Main   string `yaml:"mainURI"`
	Upload string `yaml:"uploadPart"`
}

type GrpcClient struct {
	RouterService             string `yaml:"routerService"`
	DocService                string `yaml:"docService"`
	StorageService            string `yaml:"storageService"`
	GatewayIntegrationService string `yaml:"gatewayIntegrationService"`
}

type Cors struct {
	AllowedOrigins     []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods     []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders     []string `json:"allowed_headers" yaml:"allowed_headers"`
	ExposedHeaders     []string `json:"exposed_headers" yaml:"exposed_headers"`
	AllowCredentials   bool     `json:"allow_credentials" yaml:"allow_credentials"`
	OptionsPassthrough bool     `json:"options_passthrough" yaml:"options_passthrough"`
	Debug              bool     `json:"debug" yaml:"debug"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configurations")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
