package auth

import (
	"github.com/BIQ-Cat/easyserver/db"
	"github.com/BIQ-Cat/easyserver/middlewares"
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/controllers"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	"github.com/BIQ-Cat/easyserver/routes"
)

func init() {
	db.ModelsList = append(db.ModelsList, models.Account{})
	middlewares.Middlewares = append(middlewares.Middlewares, app.JWTAuthentication)
	routes.Routes["auth"] = &controllers.Route
}
