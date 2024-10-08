package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	// Modules

	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/BIQ-Cat/easyserver/internal/utils"
	// Configuration
)

func init() {
	sendOTP := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(app.UserKey{}).(uint)
		var acc models.Account

		err := db.GetDB().Table("accounts").Where("id = ?", id).First(&acc).Error
		if err == gorm.ErrRecordNotFound {
			panic(err) // must exist
		} else if err != nil {
			utils.HandleError(w, err)
			return
		}

		var resp map[string]interface{}

		if r.URL.Query().Has("email") && moduleconfig.Config.Create.Email.Require {
			resp, err = acc.SendEmailOTP(r.URL.Query().Get("email"), false, r.Host)
		} else {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		if err != nil {
			utils.HandleError(w, err)
			return
		}
		utils.Respond(w, resp)
	}

	Route["forgot-password"] = routes.Controller{
		Handler:     http.HandlerFunc(sendOTP),
		Methods:     []string{"GET"},
		RequireAuth: true,
	}
}
