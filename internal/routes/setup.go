package routes

import (
	"fmt"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/gorilla/mux"
)

func Setup(root *mux.Router) {
	for name, subroutes := range Routes {
		if config.Config.Debug {
			fmt.Println(name, "contains")
		}

		subrouter := root.PathPrefix("/" + name).Subrouter()
		for subpath, handler := range *subroutes {
			if config.Config.Debug {
				fmt.Println("\t", subpath)
			}

			route := subrouter.Handle("/"+subpath, handler)
			if len(handler.Methods) != 0 {
				route.Methods(handler.Methods...)
			}
			if len(handler.Schemas) != 0 {
				route.Schemes(handler.Schemas...)
			}
			if len(handler.Headers) != 0 {
				route.Methods(handler.Headers...)
			}
		}
	}
}
