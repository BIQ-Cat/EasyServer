package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/BIQ-Cat/easyserver/internal/utils"
)

func init() {
	visualReset := func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("visual") || !r.URL.Query().Has("token") {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		templ, err := utils.ParseTemplateDir("auth")
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		err = templ.ExecuteTemplate(w, "change-password.html", struct {
			Token string
		}{
			Token: r.URL.Query().Get("token"),
		})

		if err != nil {
			utils.HandleError(w, err)
			return
		}
	}

	type ResetPasswordInput struct {
		Token    string
		Password string
	}

	reset := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			visualReset(w, r)
			return
		}

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var params ResetPasswordInput

		err := decoder.Decode(&params)
		if err != nil {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}
		otp, err := base64.StdEncoding.DecodeString(params.Token)
		if err != nil {
			utils.HandleError(w, err)
		}

		resp, err := models.ResetPassword(otp, []byte(params.Password))
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		utils.Respond(w, resp)
	}

	Route["reset-password"] = routes.Controller{
		Handler:     http.HandlerFunc(reset),
		Methods:     []string{http.MethodGet, http.MethodPatch},
		RequireAuth: false,
	}
}
