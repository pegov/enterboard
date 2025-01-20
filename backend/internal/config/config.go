package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	App App
}

type App struct {
	Host string `envconfig:"APP_HOST"`
	Port uint16 `envconfig:"APP_PORT"`
}

func New() *Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return &cfg
}
