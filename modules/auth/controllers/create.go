package controllers

import (
	"encoding/json"
	"net/http"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/BIQ-Cat/easyserver/internal/utils"
)

func init() {
	createUser := func(w http.ResponseWriter, r *http.Request) {
		account := &models.Account{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(account); err != nil {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		resp, err := account.Create()
		if err != nil {
			utils.HandleError(w, err)
			return
		}
		utils.Respond(w, resp)
	}

	Route["create"] = routes.Controller{
		Handler: http.HandlerFunc(createUser),
		Methods: []string{"POST"},
	}
}
