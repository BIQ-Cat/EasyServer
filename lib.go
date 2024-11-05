package easyserver

import (
	"net/http"
	"strings"
)

type Router struct {
	Modules map[string]Module
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
	Headers []string
	Schemas []string
	Data    map[string]any
}

type MiddlewareFunc func(http.Handler) http.Handler

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for path, module := range router.Modules {
		if strings.HasPrefix(r.URL.Path, path) {
			for range module.Middlewares {

			}
		}
	}
}
