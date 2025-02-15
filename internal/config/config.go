package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Env    string `yaml:"env" env:"ENV" env-default:"local" env-description:"app environment"`
	DB     `yaml:"db" env-description:"db environment"`
	Server `yaml:"server" env-description:"server environment"`
}

type DB struct {
	Host     string `yaml:"host" env-default:"avito-shop-db" env-required:"true"`
	Port     string `yaml:"port" env-default:"5432" env-required:"true"`
	Username string `yaml:"username" env-default:"postgres" env-required:"true"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name" env-default:"shop" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"disable" env-required:"true"`
}

type Server struct {
	Address        string        `yaml:"address" env-default:":8080"`
	Timeout        time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
	MaxHeaderBytes int           `yaml:"max_header_bytes:" env-default:"1048576"`
}

func Load() *Config {
	configPath := "./config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("cannot read config", "path", configPath, "error", err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		slog.Error("cannot read config", "error", err)
	}

	return &config
}
