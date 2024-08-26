package config

import "github.com/BIQ-Cat/easyserver/config/auto"

func LoadEnv() (err error) {
	EnvConfig, err = auto.ParseEnv(Config.Debug, &EnvConfig)
	return
}
