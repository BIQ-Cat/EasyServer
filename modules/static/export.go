package static

import (
	"net/http"

	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/internal/router"
)

func init() {
	router.DefaultRouter.Modules["static"] = easyserver.Module{
		Route: easyserver.Route{
			"/": easyserver.Controller{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.ServeFile(w, r, "./static/index.html")
				}),
			},
			"/singin": easyserver.Controller{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.ServeFile(w, r, "./static/singin.html")
				}),
			},
		},
	}

	// var cfg basictypes.JSONConfig = moduleconfig.Config
	// json.Configurations["static"] = &cfg
}
