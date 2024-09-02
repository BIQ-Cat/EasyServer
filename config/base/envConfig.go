package config

import (
	"github.com/BIQ-Cat/easyserver/config/base/funcs"
	"github.com/BIQ-Cat/easyserver/config/base/types"
)

var EnvConfig = types.EnvConfig{
	TokenPassword: funcs.GenerateTokenPassword(30),
	OTPPassword:   funcs.GenerateTokenPassword(30),
}
