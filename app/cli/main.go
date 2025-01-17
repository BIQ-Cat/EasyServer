package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BIQ-Cat/easyserver/internal/json"
	"github.com/BIQ-Cat/easyserver/internal/router"

	// Module imports
	_ "github.com/BIQ-Cat/easyserver/config/modules"

	// Configuration
	config "github.com/BIQ-Cat/easyserver/config/base"
	configFuncs "github.com/BIQ-Cat/easyserver/config/base/funcs"
)

func main() {
	loadEnv()

	err := router.DefaultRouter.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = json.ParseFiles()
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      &router.DefaultRouter,
		Addr:         fmt.Sprintf(":%d", config.EnvConfig.ServerPort),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server is running on port", config.EnvConfig.ServerPort)
	log.Fatal(srv.ListenAndServe())
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
