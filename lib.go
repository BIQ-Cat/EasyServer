package easyserver

import (
	"net/http"
	"path"
	"strings"
)

type Router struct {
	Modules map[string]Module
}

type Module struct {
	Route
	Middlewares    []MiddlewareFunc
	HTTP404Handler http.HandlerFunc
	EnableAsterix  bool
	Models         []interface{}
}

type Route map[string]Controller

type Controller struct {
	http.Handler
	Methods []string
	Headers []string
	Schemas []string
	Data    map[string]any
}

type MiddlewareFunc http.HandlerFunc

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for modulePath, module := range router.Modules {
		relativePath, hasPrefix := strings.CutPrefix(r.URL.Path, "/"+modulePath)
		if hasPrefix {
			controllerPath := strings.TrimPrefix(relativePath, "/")

			var handler http.Handler
			handler, ok := module.Route[controllerPath]
			if !ok {
				if controllerPath == "" {
					w.WriteHeader(404)
					module.HTTP404Handler.ServeHTTP(w, r)
					return
				}

				parentPath := path.Dir(controllerPath) + "/"
				if parentPath == "//" {
					parentPath = ""
				}

				handler, ok = module.Route[parentPath+"..."]

				if !ok && module.EnableAsterix {
					if parentPath == "" {
						handler, ok = module.Route["*"]
						if !ok {

							module.HTTP404Handler.ServeHTTP(w, r)
							return
						}
					}
				}
			}

			for range module.Middlewares {
				handler.ServeHTTP(w, r)
			}
		}
	}
}
