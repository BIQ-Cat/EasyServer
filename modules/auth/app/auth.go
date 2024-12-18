package app

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	// Modules
	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/modules/auth/datakeys"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/router"
	"github.com/BIQ-Cat/easyserver/internal/utils"

	// Configuration
	config "github.com/BIQ-Cat/easyserver/config/base"
	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
)

type UserKey struct{}

func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := router.DefaultRouter.GetController(r.URL.Path)
		if _, ok := controller.Data[datakeys.RequireAuth]; controller == nil || !ok {
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

		bearerToken := strings.Split(tokenHeader, " ") // []string{"Bearer", "<token>"} (?)
		if len(bearerToken) != 2 && bearerToken[0] != "Bearer" {
			response = utils.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			utils.Respond(w, response)
			return
		}

		tokenPart := bearerToken[1]
		tk := &models.Token{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(_ *jwt.Token) (interface{}, error) {
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
			fmt.Printf("User %v\n", tk.UserID)
		}

		if verificationRequired(*tk, r.URL.Path, controller) {
			response = utils.Message(false, "Verification required.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			utils.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func verificationRequired(tk models.Token, path string, c *easyserver.Controller) bool {
	_, ok := c.Data[datakeys.RequireVerification]
	return !tk.Verified && (moduleconfig.Config.Verify.Require && path != "/auth/verify-send" && path != "/auth/verify-recieve" || ok)
}
