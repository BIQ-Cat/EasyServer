package auth

import (
	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/controllers"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/json"
	"github.com/BIQ-Cat/easyserver/internal/middlewares"
	"github.com/BIQ-Cat/easyserver/internal/routes"

	// Configuration
	moduleConfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	basicTypes "github.com/BIQ-Cat/easyserver/config/types"
)

func init() {
	db.ModelsList = append(db.ModelsList, models.Account{})
	middlewares.Middlewares = append(middlewares.Middlewares, app.JWTAuthentication)
	routes.Routes["auth"] = &controllers.Route

	var cfg basicTypes.JSONConfig = moduleConfig.Config
	json.Configurations = append(json.Configurations, &cfg)
}
