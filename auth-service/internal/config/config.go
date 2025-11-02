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
		Host     string `env:"MONGODB_HOST" env-default:"mongo"`
		User     string `env:"MONGODB_USER" env-default:"user"`
		Pass     string `env:"MONGODB_PASS" env-default:"pass"`
		Database string `env:"MONGODB_DATABASE" env-default:"test"`
	}
	App struct {
		ENV         string `env:"APP_ENV" env-default:"development"`
		StorageType string `env:"APP_STORAGE" env-default:"mongodb"`
		SecretKey   string `env:"SECRET_KEY" env-default:"my_secret_key"`
	}
}

var cfg Config

func MustLoad() *Config {
	if cfg == (Config{}) {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(fmt.Errorf("problem with read env %s", err))
		}
	}

	return &cfg
}
