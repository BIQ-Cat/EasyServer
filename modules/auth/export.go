package auth

import (
	// Modules
	"github.com/BIQ-Cat/easyserver"

	// Internals

	"github.com/BIQ-Cat/easyserver/internal/router"
	// Configuration
)

func init() {
	router.DefaultRouter.Modules["auth"] = easyserver.Module{
		// Route:       controllers.Route,
		// Middlewares: []easyserver.MiddlewareFunc{app.JWTAuthentication},
		// Models:      []any{models.Account{}},
	}

	// var cfg basictypes.JSONConfig = moduleconfig.Config
	// json.Configurations["auth"] = &cfg
}
