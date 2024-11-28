package auth

import (
	"github.com/gorilla/mux"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/controllers"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/addons"

	// Configuration
	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	basictypes "github.com/BIQ-Cat/easyserver/config/types"
)

var cfg basictypes.JSONConfig = moduleconfig.Config

var Module = addons.Module{
	Middlewares:   []mux.MiddlewareFunc{app.JWTAuthentication},
	Models:        []interface{}{models.Account{}},
	Route:         &controllers.Route,
	Configuration: &cfg,
}

func init() {
	ok := Module.Register("auth", "auth")
	if !ok {
		panic("module auth already registered")
	}
}
