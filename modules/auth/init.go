package auth

import (
	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/controllers"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/middlewares"
	"github.com/BIQ-Cat/easyserver/internal/routes"
)

func init() {
	db.ModelsList = append(db.ModelsList, models.Account{})
	middlewares.Middlewares = append(middlewares.Middlewares, app.JWTAuthentication)
	routes.Routes["auth"] = &controllers.Route
}
