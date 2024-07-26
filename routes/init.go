package routes

import (
	"net/http"
)

type Route struct {
	http.Handler
	Methods     []string
	Headers     []string
	Schemas     []string
	RequireAuth bool
}

var Routes map[string]*map[string]Route
