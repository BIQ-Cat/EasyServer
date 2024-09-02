package routes

import (
	"net/http"
)

type Controller struct {
	http.Handler
	Methods             []string
	Headers             []string
	Schemas             []string
	RequireAuth         bool
	RequireVerification bool
}

type Route map[string]Controller

var Routes = make(map[string]*Route)
