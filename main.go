package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BIQ-Cat/easyserver/app"
	"github.com/BIQ-Cat/easyserver/routes"
	_ "github.com/BIQ-Cat/easyserver/routes/auth"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var debug = flag.Bool("debug", false, "Enable debug logs")

func main() {
	flag.Parse()
	utils.SetDebug(*debug)

	root := mux.NewRouter()
	root.Use(app.JWTAuthentication)

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
