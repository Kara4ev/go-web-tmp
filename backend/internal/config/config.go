package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      `yaml:"app"`
		HTTP     `yaml:"http"`
		Logger   `yaml:"logger"`
		Postgres `yaml:"postgres"`
		Redis    `yaml:"radis"`
	}

	App struct {
		AppName                   string `yaml:"name" env-required:"true"`
		AppVersion                string `yaml:"version" env-required:"true"`
		AppDebug                  int    `yaml:"debug" env-required:"true" env:"APP_DEBUG"`
		AppSecret                 string `yaml:"secret" env-required:"true" env:"APP_SECRET"`
		AppPrivateKeyFile         string `yaml:"privat_key_file" env-required:"true" env:"APP_PRIV_KEY_FILE"`
		AppPublicKeyFile          string `yaml:"pub_key_file" env-required:"true" env:"APP_PUB_KEY_FILE"`
		AppLogFile                string `yaml:"log_file" env-required:"true" env:"APP_LOG_FILE"`
		AppRefreshTokenExpiration int64  `yaml:"refresh_token_exp" env-required:"true" env:"APP_R_TOKEN_EXP"`
		AppIDTokenExpiration      int64  `yaml:"id_token_exp" env-required:"true" env:"APP_R_TOKEN_EXP"`
	}

	HTTP struct {
		HTTPHost           string `yaml:"host" env-required:"true" env:"HTTP_HOST"`
		HTTPPort           string `yaml:"port" env-required:"true" env:"HTTP_PORT"`
		HTTPBaseURL        string `yaml:"base_url" env-required:"true" env:"HTTP_BASE_URL"`
		HTTPHendlerTimeOut int64  `yaml:"hendler_time_out" env-required:"true" env:"HTTP_HENDLER_TIME_OUT"`
	}

	Logger struct {
		LoggerLevel string `yaml:"level" env-required:"true" env:"LOG_LEVEL"`
	}

	Postgres struct {
		PGHost     string `yaml:"host" env-required:"true" env:"PG_HOST"`
		PGPort     string `yaml:"port" env-required:"true" env:"PG_PORT"`
		PGUser     string `yaml:"user" env-required:"true" env:"PG_USER"`
		PGPassword string `yaml:"password" env-required:"true" env:"PG_PASSWORD"`
		PGDB       string `yaml:"db" env-required:"true" env:"PG_DB"`
		PGSSL      string `yaml:"ssl" env-required:"true" env:"PG_SSL"`
	}

	Redis struct {
		RDHost     string `yaml:"host" env-required:"true" env:"RD_HOST"`
		RDPort     string `yaml:"port" env-required:"true" env:"RD_PORT"`
		RDPassword string `yaml:"password" env-required:"true" env:"RD_PASSWORD"`
		RDdb       int    `yaml:"db" env-required:"true" env:"RD_DB"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	pathFileConfig := "./config/config.yaml"

	if err := cleanenv.ReadConfig(pathFileConfig, cfg); err != nil {
		return nil, fmt.Errorf("error init file config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error init env config: %w", err)
	}

	return cfg, nil

}
