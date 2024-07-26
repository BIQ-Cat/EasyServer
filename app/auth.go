package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/BIQ-Cat/easyserver/models"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/dgrijalva/jwt-go"
)

type UserKey struct{}

func checkNotAuth(requestPath string) bool {
	var wg sync.WaitGroup
	var doNotAuth chan bool
	for path, subroutes := range routes.Routes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for subpath, route := range *subroutes {
				if !route.RequireAuth {
					fullPath := fmt.Sprintf("/%s/%s", path, subpath)
					if fullPath == requestPath {
						doNotAuth <- true
					}
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		doNotAuth <- false
	}()
	return <-doNotAuth
}

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkNotAuth(r.URL.Path) {
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
		tk := models.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
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

		utils.LogInDebug(fmt.Sprintf("User %v", tk.UserId))
		ctx := context.WithValue(r.Context(), UserKey{}, tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
