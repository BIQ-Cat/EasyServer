package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"

	// Modules

	moduleconfig "github.com/BIQ-Cat/easyserver/config/modules/auth"
	"github.com/BIQ-Cat/easyserver/modules/auth/app"
	"github.com/BIQ-Cat/easyserver/modules/auth/datakeys"
	"github.com/BIQ-Cat/easyserver/modules/auth/models"

	// Internals
	"github.com/BIQ-Cat/easyserver/internal/db"
	"github.com/BIQ-Cat/easyserver/internal/routes"
	"github.com/BIQ-Cat/easyserver/internal/utils"
	// Configuration
)

func init() {
	changePassword := func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(app.UserKey{}).(uint)
		var acc models.Account

		err := db.GetDB().Table("accounts").Where("id = ?", id).First(&acc).Error
		if err == gorm.ErrRecordNotFound {
			panic(err) // must exist
		} else if err != nil {
			utils.HandleError(w, err)
			return
		}

		if !moduleconfig.Config.Verify.SetPasswordAfter || acc.Password != "" {
			utils.Respond(w, utils.Message(false, "Forbidden"))
			return
		}

		err = r.ParseForm()
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		if r.Form.Has("password") {
			utils.Respond(w, utils.Message(false, "Invalid request"))
			return
		}

		resp, err := acc.ChangePassword([]byte(r.Form.Get("password")))
		if err != nil {
			utils.HandleError(w, err)
		}
		utils.Respond(w, resp)
	}

	Route["change-password"] = routes.Controller{
		Handler: http.HandlerFunc(changePassword),
		Methods: []string{http.MethodPost},
		Data: map[string]any{
			datakeys.RequireAuth: true,
		},
	}
}
