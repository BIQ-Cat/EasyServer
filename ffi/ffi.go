package main

import "C"

import (
	"encoding/json"
	"reflect"
	"unsafe"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/BIQ-Cat/easyserver/config/base/funcs"
	"github.com/BIQ-Cat/easyserver/internal/addons"
	"github.com/BIQ-Cat/easyserver/internal/db"

	_ "github.com/BIQ-Cat/easyserver/config/modules"
)

//export GetDefaultModuleConfiguration
func GetDefaultModuleConfiguration(moduleName string) (data *C.char, ok bool) {
	var module *addons.Module
	module, ok = addons.GetModule(moduleName)
	if !ok || module.Configuration == nil {
		ok = false
		return
	}

	res, err := json.MarshalIndent(module.Configuration, "", "  ")
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
	cArray := C.malloc(C.size_t(len(db.ModelsList)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	goSlice := unsafe.Slice((**C.char)(cArray), len(db.ModelsList))
	for i, model := range db.ModelsList {
		goSlice[i] = C.CString(reflect.TypeOf(model).Name())
	}

	return (C.size_t)(len(db.ModelsList)), (**C.char)(cArray)
}

//export ListModules
func ListModules() (C.size_t, **C.char) {
	moduleNames := addons.GetModuleNames()

	cArray := C.malloc(C.size_t(len(moduleNames)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	goSlice := unsafe.Slice((**C.char)(cArray), len(moduleNames))
	for i, module := range moduleNames {
		goSlice[i] = C.CString(module)
	}

	return (C.size_t)(len(moduleNames)), (**C.char)(cArray)
}

func main() {}
