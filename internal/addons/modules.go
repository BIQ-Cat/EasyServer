package addons

import (
	basictypes "github.com/BIQ-Cat/easyserver/config/types"
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/json"
	"github.com/BIQ-Cat/easyserver/internal/middlewares"
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/gorilla/mux"
)

type Module struct {
	Route         *routes.Route
	Models        []interface{}
	Configuration *basictypes.JSONConfig
	Middlewares   []mux.MiddlewareFunc
}

var modules map[string]*Module

func init() {
	if modules == nil {
		modules = make(map[string]*Module)
	}
}

func (m *Module) Register(moduleName string, routeName string) (ok bool) {
	_, ok = modules[moduleName]
	ok = !ok
	if !ok {
		return
	}

	modules[moduleName] = m

	if m.Models != nil {
		db.ModelsList = append(db.ModelsList, m.Models...)
	}

	if m.Route != nil {
		routes.Routes[routeName] = m.Route
	}

	if m.Configuration != nil {
		json.Configurations = append(json.Configurations, m.Configuration)
	}

	if m.Middlewares != nil {
		middlewares.Middlewares = append(middlewares.Middlewares, m.Middlewares...)
	}
	return
}

func GetModule(moduleName string) (module *Module, ok bool) {
	module, ok = modules[moduleName]
	return
}

func GetModuleNames() []string {
	keys := make([]string, 0, len(modules))
	for key := range modules {
		keys = append(keys, key)
	}
	return keys
}
