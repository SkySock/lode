package app

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env" env-default:"local"`
	HTTP HTTPConfig `yaml:"http"`
	Auth AuthConfig `yaml:"auth"`
}

type HTTPConfig struct {
	Host string `yaml:"host" env-default:"127.0.0.1"`
	Port uint16 `yaml:"port" env-default:"8000"`
}

type AuthConfig struct {
	AccessSecretKey string `yaml:"access_secret"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
