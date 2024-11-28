package core

import (
	"github.com/ilyakaznacheev/cleanenv"
	"net"
	"strconv"
)

type Config struct {
	Host           string   `json:"host" env-default:"0.0.0.0" env:"HOST"`
	Port           int      `json:"port" env-default:"8080" env:"PORT"`
	AllowedOrigins []string `json:"allowed_origins" env:"ALLOWED_ORIGINS" default:"*"`
	AllowedHeaders []string `json:"allowed_headers" env:"ALLOWED_HEADERS" default:"*"`
	AllowedMethods []string `json:"allowed_methods" env:"ALLOWED_METHODS" default:"GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD"`
}

func (cfg *Config) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Config Config `json:"core" env-prefix:"HTTP_"`
	}
	err := cleanenv.ReadConfig("config.json", &cfg)
	if err != nil {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}
	return &cfg.Config, nil
}
