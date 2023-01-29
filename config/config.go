package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
)

var Cfg Config

type Config struct {
	Debug bool `env:"DEBUG" env-default:"False"`
	MySQL struct {
		Port string `env:"MYSQL_PORT" env-default:"3306"`
		Host string `env:"MYSQL_HOST" env-default:"pdns-mysql"`
		Name string `env:"MYSQL_NAME" env-default:"pdns"`
		User string `env:"MYSQL_USER" env-default:"root"`
		Pass string `env:"MYSQL_PASS" env-default:"root"`
	}
	PowerDNS struct {
		Host     string `env:"PDNS_HOST" env-default:"pdns:8081"`
		XAPIKey  string `env:"PDNS_X-API-Key" env-default:"gngjngogjrngoengong"`
		ServerID string `env:"PDNS_SERVERID" env-default:"pdns-mysql"`
	}
	NS struct {
		NS1 string `env:"NS1" env-default:"NS1.ABR.OOO"`
		NS2 string `env:"NS2" env-default:"NS2.ABR.OOO"`
	}
	Log struct {
		LogLevel   string `env:"LOGLEVEL" env-default:"debug"`
		OutputType string `env:"LOG_OUTPUT_TYPE" env-default:"stdout"`
		OutputAdd  string `env:"LOG_FILE_Add" env-default:"/log.txt"`
	}
	Auth struct {
		Host string `env:"AUTH_HOST" env-default:"https://api.abr.ooo/v0/user"`
	}
	Sentry struct {
		DSN   string `env:"SENTRY_DSN" env-default:"sentry_dsn_address"`
		Level string `env:"SENTRY_LEVEL" env-default:"error"`
	}
}

func SetConfig() {

	if _, err := os.Stat(".env"); err == nil {
		cleanenv.ReadConfig(".env", &Cfg)
		logrus.Info("Set config from .env file")
	} else {
		cleanenv.ReadEnv(&Cfg)
		logrus.Info("Set config from Config struct values")
	}

}
