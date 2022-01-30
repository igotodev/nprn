package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"nprn/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug bool    `yaml:"is_debug"`
	Listen  Listen  `yaml:"listen"`
	MongoDB MongoDB `yaml:"mongo_db"`
}

type Listen struct {
	Type   string `yaml:"type" env-default:"tcp"`
	Port   string `yaml:"port" env-default:"8080"`
	BindIP string `yaml:"bind_ip" env-default:"0.0.0.0"`
}

type MongoDB struct {
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	DBName         string `yaml:"db_name"`
	UserCollection string `yaml:"user_collection"`
	SaleCollection string `yaml:"sale_collection"`
	AuthDB         string `yaml:"auth_db"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
}

var instance *Config
var once sync.Once

// singleton
func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")

		instance = &Config{}

		err := cleanenv.ReadConfig("config.yaml", instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}

	})
	return instance
}
