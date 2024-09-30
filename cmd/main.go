package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/json"
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

	err = json.ParseFiles()
	if err != nil {
		log.Fatal(err)
	}

	root := mux.NewRouter()
	root.Use(middlewares.Middlewares...)

	routes.Setup(root)

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
