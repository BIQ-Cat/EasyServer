package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BIQ-Cat/easyserver/config"
	"github.com/BIQ-Cat/easyserver/config/auto"
	"github.com/BIQ-Cat/easyserver/db"
	"github.com/BIQ-Cat/easyserver/middlewares"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/gorilla/mux"

	// Module imports
	_ "github.com/BIQ-Cat/easyserver/modules"
)

func main() {
	err := config.LoadEnv()
	if err == auto.ErrEnvNotSet {
		log.Fatal("Not all environment variables are set")
	} else if err != nil {
		log.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	root := mux.NewRouter()
	root.Use(middlewares.Middlewares...)

	for name, subroutes := range routes.Routes {
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

	fmt.Println("Server is running on port", config.EnvConfig.ServerPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.EnvConfig.ServerPort), root))
}
