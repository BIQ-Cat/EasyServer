package controllers

import (
	"net/http"

	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/BIQ-Cat/easyserver/utils"
)

func init() {
	verifyOTP := func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("token") {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		resp, err := models.VerifyAccount(r.URL.Query().Get("token"))
		if err != nil {
			utils.HandleError(w, err)
			return
		}
		if r.URL.Query().Has("visual") {
			type param struct {
				Status bool
			}

			templ, err := utils.ParseTemplateDir("templates")
			if err != nil {
				utils.HandleError(w, err)
				return
			}

			err = templ.ExecuteTemplate(w, "verify-recieve.html", &param{utils.GetStatus(resp)})
			if err != nil {
				utils.HandleError(w, err)
				return
			}
		}
		utils.Respond(w, resp)
	}

	Route["verify-recieve"] = routes.Controller{
		Handler:     http.HandlerFunc(verifyOTP),
		Methods:     []string{"GET"},
		RequireAuth: false,
	}
}
