package controllers

import (
	"encoding/base64"
	"net/http"

	// Modules
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/BIQ-Cat/easyserver/internal/utils"
)

func init() {
	respond := func(w http.ResponseWriter, r *http.Request, resp map[string]interface{}) {
		if r.URL.Query().Has("visual") {
			templ, err := utils.ParseTemplateDir("templates")
			if err != nil {
				utils.HandleError(w, err)
				return
			}

			err = templ.ExecuteTemplate(w, "verify-recieve.html", &resp)
			if err != nil {
				utils.HandleError(w, err)
				return
			}
		} else {
			utils.Respond(w, resp)
		}
	}

	verifyOTP := func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("token") {
			respond(w, r, utils.Message(false, "Invalid request"))
			return
		}

		token := r.URL.Query().Get("token")
		otp, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			utils.HandleError(w, err)
		}

		resp, err := models.VerifyAccount(otp)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		respond(w, r, resp)
	}

	Route["verify-recieve"] = routes.Controller{
		Handler: http.HandlerFunc(verifyOTP),
		Methods: []string{"GET"},
	}
}
