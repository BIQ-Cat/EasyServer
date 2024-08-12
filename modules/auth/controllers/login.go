package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BIQ-Cat/easyserver/config"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/BIQ-Cat/easyserver/utils"
)

func init() {
	logIn := func(w http.ResponseWriter, r *http.Request) {
		account := &models.Account{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(account); err != nil {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		resp, err := models.Login(account.Username, account.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if config.Config.Debug {
				fmt.Fprint(w, err)
			}
		}
		utils.Respond(w, resp)
	}

	Route["login"] = routes.Controller{
		Handler: http.HandlerFunc(logIn),
		Methods: []string{"POST"},
	}
}
