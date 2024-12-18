package main

import "C"

import (
	"encoding/json"
	"reflect"
	"unsafe"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/BIQ-Cat/easyserver/config/base/funcs"
	basictypes "github.com/BIQ-Cat/easyserver/config/types"
	jsonConfig "github.com/BIQ-Cat/easyserver/internal/json"
	"github.com/BIQ-Cat/easyserver/internal/router"

	_ "github.com/BIQ-Cat/easyserver/config/modules"
)

//export GetDefaultModuleConfiguration
func GetDefaultModuleConfiguration(moduleName string) (data *C.char, ok bool) {
	var config *basictypes.JSONConfig
	config, ok = jsonConfig.Configurations[moduleName]
	if !ok || config == nil {
		ok = false
		return
	}

	res, err := json.MarshalIndent(config, "", "  ")
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

//export ListModels
func ListModels() (C.size_t, **C.char) {
	modelsList := router.DefaultRouter.ModelsList()
	cArray := C.malloc(C.size_t(len(modelsList)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	goSlice := unsafe.Slice((**C.char)(cArray), len(modelsList))
	for i, model := range modelsList {
		goSlice[i] = C.CString(reflect.TypeOf(model).Name())
	}

	return (C.size_t)(len(modelsList)), (**C.char)(cArray)
}

//export ListModules
func ListModules() (C.size_t, **C.char) {
	moduleNames := router.DefaultRouter.ModuleNames()

	cArray := C.malloc(C.size_t(len(moduleNames)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	goSlice := unsafe.Slice((**C.char)(cArray), len(moduleNames))
	for i, module := range moduleNames {
		goSlice[i] = C.CString(module)
	}

	return (C.size_t)(len(moduleNames)), (**C.char)(cArray)
}

func main() {}
