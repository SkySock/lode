package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string       `yaml:"env" env-default:"local"`
	HTTP   HTTPConfig   `yaml:"http"`
	DB     DBConfig     `yaml:"db"`
	Valkey ValkeyConfig `yaml:"valkey"`
	Auth   AuthConfig   `yaml:"auth"`
}

type HTTPConfig struct {
	Host string `yaml:"host" env-default:"127.0.0.1"`
	Port uint16 `yaml:"port" env-default:"8000"`
}

type DBConfig struct {
	URL     string `yaml:"url"`
	Migrate bool   `yaml:"migrate" env-default:"false"`
}

type ValkeyConfig struct {
	Addr string `yaml:"addr" env-default:"127.0.0.1:6379"`
}

type AuthConfig struct {
	AccessSecretKey string        `yaml:"access_secret"`
	Lifetime        TokenLifetime `yaml:"lifetime"`
}

type TokenLifetime struct {
	Access  time.Duration `yaml:"access" env-default:"5m"`
	Refresh time.Duration `yaml:"refresh" env-default:"720h"`
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
