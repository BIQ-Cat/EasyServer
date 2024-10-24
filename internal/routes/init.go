package routes

import (
	"net/http"
)

type Controller struct {
	http.Handler
	Methods []string
	Headers []string
	Schemas []string
	Data    map[string]any
}

type Route map[string]Controller

var Routes = make(map[string]*Route)
