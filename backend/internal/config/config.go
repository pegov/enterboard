package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App   App
	DB    DB
	Cache Cache
}

type App struct {
	Host string `envconfig:"APP_HOST"`
	Port uint16 `envconfig:"APP_PORT"`
}

type DB struct {
	Host            string        `envconfig:"DB_HOST"`
	Port            uint16        `envconfig:"DB_PORT"`
	Database        string        `envconfig:"DB_DATABASE"`
	Username        string        `envconfig:"DB_USERNAME"`
	Password        string        `envconfig:"DB_PASSWORD"`
	MaxIdleConns    int           `envconfig:"DB_MAX_IDLE_CONNS" default:"20"`
	MaxOpenConns    int           `envconfig:"DB_MAX_OPEN_CONNS" default:"20"`
	ConnMaxLifetime time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"1h"`
}

type Cache struct {
	Host     string `envconfig:"CACHE_HOST"`
	Port     uint16 `envconfig:"CACHE_PORT"`
	Database string `envconfig:"CACHE_DATABASE"`
}

func New() *Config {
	var cfg Config
	envconfig.MustProcess("", &cfg)
	return &cfg
}
