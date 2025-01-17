package auth

import (
	// Modules
	"github.com/BIQ-Cat/easyserver"

	// Internals

	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	basictypes "github.com/BIQ-Cat/easyserver/config/types"
	"github.com/BIQ-Cat/easyserver/internal/json"
	"github.com/BIQ-Cat/easyserver/internal/router"
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/controllers"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	// Configuration
)

func init() {
	router.DefaultRouter.Modules["auth"] = easyserver.Module{
		Route:       controllers.Route,
		Middlewares: []easyserver.MiddlewareFunc{app.JWTAuthentication},
		Models:      []any{models.Account{}},
	}

	var cfg basictypes.JSONConfig = moduleconfig.Config
	json.Configurations["auth"] = &cfg
}
