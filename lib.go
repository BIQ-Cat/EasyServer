package easyserver

import (
	"net/http"
)

type Router struct {
	http.ServeMux
}

func NewRouter() *Router {
	return new(Router)
}