package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/middlewares"
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/gorilla/mux"

	// Module imports
	_ "github.com/BIQ-Cat/easyserver/config/modules"

	// Configuration
	config "github.com/BIQ-Cat/easyserver/config/base"
	configFuncs "github.com/BIQ-Cat/easyserver/config/base/funcs"
)

func main() {
	loadEnv()

	err := db.Connect()
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

func loadEnv() {
	var err error

	config.EnvConfig, err = configFuncs.ParseEnv(config.Config.Debug, &config.EnvConfig)
	if err == configFuncs.ErrEnvNotSet {
		log.Fatal("Not all environment variables are set")
	} else if err != nil {
		log.Fatal(err)
	}
}
