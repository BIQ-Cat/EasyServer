package main

import "C"

import (
	"encoding/json"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/BIQ-Cat/easyserver/config/base/funcs"
	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	basictypes "github.com/BIQ-Cat/easyserver/config/types"
)

var Configs = map[string]basictypes.JSONConfig{
	"auth": moduleconfig.Config,
}

//export GetDefaultModuleConfiguration
func GetDefaultModuleConfiguration(moduleName string) (data *C.char, ok bool) {
	var cfg basictypes.JSONConfig
	cfg, ok = Configs[moduleName]
	if !ok || !cfg.HasExternalFile() {
		ok = false
		return
	}

	res, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		ok = false
	}
	data = C.CString(string(res))
	return
}

//export GetEnvironmentConfiguration
func GetEnvironmentConfiguration() (data *C.char, ok bool) {
	env := config.EnvConfig
	defaults, err := funcs.EnvConfigToMap(&env)
	if err != nil {
		return
	}

	goData, err := json.Marshal(defaults)
	if err != nil {
		return
	}
	ok = true
	data = C.CString(string(goData))
	return
}

func main() {}
