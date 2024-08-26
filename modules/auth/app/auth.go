package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/BIQ-Cat/easyserver/config"
	moduleConfig "github.com/BIQ-Cat/easyserver/modules/auth/config"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/golang-jwt/jwt/v5"
)

type UserKey struct{}

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := getController(r.URL.Path)
		if controller == nil || !controller.RequireAuth {
			next.ServeHTTP(w, r)
			return
		}
		var response map[string]interface{}
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			response = utils.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, response)
			return
		}
		tokenPart := splitted[1]
		tk := &models.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.EnvConfig.TokenPassword), nil
		})

		if err != nil {
			response = utils.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, response)
			return
		}

		if !token.Valid {
			response = utils.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		if config.Config.Debug {
			fmt.Printf("User %v\n", tk.UserId)
		}

		if verificationRequired(*tk, r.URL.Path, controller) {
			response = utils.Message(false, "Verification required.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getController(requestPath string) *routes.Controller {
	for path, subroutes := range routes.Routes {
		for subpath, route := range *subroutes {
			fullPath := fmt.Sprintf("/%s/%s", path, subpath)
			if fullPath == requestPath {
				return &route
			}
		}
	}

	if config.Config.Debug {
		log.Println("WARNING: no controller found")
	}
	return nil // 404
}

func verificationRequired(tk models.Token, path string, c *routes.Controller) bool {
	return !tk.Verified && (moduleConfig.Config.Verify.Require && path != "/auth/verify-send" && path != "/auth/verify-recieve" || c.RequireVerification)
}
