package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BIQ-Cat/easyserver/db"
	"github.com/BIQ-Cat/easyserver/middlewares"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	// Module imports
	_ "github.com/BIQ-Cat/easyserver/modules"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	root := mux.NewRouter()
	root.Use(middlewares.Middlewares...)

	for name, subroutes := range routes.Routes {
		subrouter := root.PathPrefix("/" + name).Subrouter()
		for subpath, handler := range *subroutes {
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, root))
}
