package controllers

import (
	"net/http"

	"github.com/BIQ-Cat/easyserver/db"
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"
	"github.com/BIQ-Cat/easyserver/routes"
	"github.com/BIQ-Cat/easyserver/utils"
	"github.com/jinzhu/gorm"
)

func init() {
	sendOTP := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(app.UserKey{}).(uint)
		var acc models.Account

		if !r.URL.Query().Has("email") {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		err := db.GetDB().Table("accounts").Where("id = ?", id).First(&acc).Error
		if err == gorm.ErrRecordNotFound {
			panic(err) // must exist
		} else if err != nil {
			utils.HandleError(w, err)
			return
		}

		resp, err := acc.SendEmailOTP(r.URL.Query().Get("email"), true, r.Host)
		if err != nil {
			utils.HandleError(w, err)
			return
		}
		utils.Respond(w, resp)
	}

	Route["verify-send"] = routes.Controller{
		Handler:     http.HandlerFunc(sendOTP),
		Methods:     []string{"GET"},
		RequireAuth: true,
	}
}
