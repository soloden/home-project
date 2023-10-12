package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer struct {
		Addr string `env:"HTTP_PORT" env-default:":9000"`
	}
	MongoDB struct {
		URL      string `env:"MONGODB_URL" env-default:"mongodb://user:pass@mongo:27017/"`
		Database string `env:"MONGODB_DATABASE" env-default:"test"`
	}
	App struct {
		StorageType string `env:"APP_STORAGE" env-default:"mongodb"`
	}
}

var cfg Config

func LoadConfig() (*Config, error) {
	if cfg == (Config{}) {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, fmt.Errorf("problem with read env %s", err)
		}
	}

	return &cfg, nil
}