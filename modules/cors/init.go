package cors

import (
	// Internals
	"github.com/BIQ-Cat/easyserver/internal/middlewares"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/cors/app"
)

func init() {
	middlewares.Middlewares = append(middlewares.Middlewares, app.EnableCORS)
}
