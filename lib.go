package easyserver

import (
	"net/http"
	pathlib "path"
	"slices"
	"strings"
)

type Router struct {
	Modules    map[string]Module
	NotExplict bool
}

type Module struct {
	Route
	Middlewares []MiddlewareFunc
	Models      []interface{}
}

type Route map[string]Controller

type Controller struct {
	http.Handler
	Methods []string
	Headers map[string]string
	Schemas []string
	Data    map[string]any
}

type MiddlewareFunc func(http.Handler) http.Handler

func (router *Router) NoMiddlewares(w http.ResponseWriter, r *http.Request) {
	newRouter := router.withoutMiddlewares()
	newRouter.ServeHTTP(w, r)
}

func (router *Router) withoutMiddlewares() http.Handler {
	newRouter := new(Router)
	for path, module := range router.Modules {
		newRouter.Modules[path] = Module{
			Route:       module.Route,
			Middlewares: []MiddlewareFunc{},
			Models:      module.Models,
		}
	}

	return newRouter
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	var gotMiddleware bool
	handler = router.withoutMiddlewares()
	for path, module := range router.Modules {
		if strings.HasPrefix(r.URL.Path, path) {
			for _, middlew := range module.Middlewares {
				gotMiddleware = true
				handler = middlew(handler)
			}
		}
	}
	if gotMiddleware {
		handler.ServeHTTP(w, r)
		return
	}

	for path, module := range router.Modules {
		if strings.HasPrefix(r.URL.Path, path) {
		loop:
			for controllerPath, controller := range module.Route {
				for header, value := range controller.Headers {
					if r.Header.Get(header) != value {
						continue loop
					}
				}
				if !slices.Contains(controller.Methods, r.Method) || !slices.Contains(controller.Schemas, r.URL.Scheme) {
					continue loop
				}
				if router.NotExplict || strings.HasSuffix(controllerPath, "*") {
					notExplictPath := strings.TrimSuffix(controllerPath, "*")
					if strings.HasPrefix(r.URL.Path, pathlib.Join(path, notExplictPath)) {
						controller.ServeHTTP(w, r)
						return
					}
				} else {
					if pathlib.Join(path, controllerPath) == r.URL.Path {
						controller.ServeHTTP(w, r)
						return
					}
				}
			}
		}
	}
}
