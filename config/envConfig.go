package config

import (
	"github.com/BIQ-Cat/easyserver/config/auto"
)

var EnvConfig = auto.EnvConfig{
	TokenPassword: auto.GenerateTokenPassword(30),
	OTPPassword:   auto.GenerateTokenPassword(30),
}
