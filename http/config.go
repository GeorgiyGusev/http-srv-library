package http

import (
	"github.com/ilyakaznacheev/cleanenv"
	"net"
	"strconv"
)

type Config struct {
	Host string `json:"host" env-default:"0.0.0.0" env:"HOST"`
	Port int    `json:"port" env-default:"8080" env:"PORT"`
}

func (cfg *Config) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}

func LoadConfig() (*Config, error) {
	var cfg struct {
		Config Config `json:"http" env-prefix:"HTTP_"`
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
