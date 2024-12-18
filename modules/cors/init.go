package cors

import (
	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/internal/router"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/cors/app"
)

func init() {
	router.DefaultRouter.Modules["cors"] = easyserver.Module{
		Middlewares: []easyserver.MiddlewareFunc{app.EnableCORS},
	}
}
