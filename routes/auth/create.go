package auth

import (
	"encoding/json"
	"net/http"

	"github.com/BIQ-Cat/easyserver/models"
	"github.com/BIQ-Cat/easyserver/routes"
)

func init() {
	createUser := func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var user models.RawUser
		decoder.Decode(&user)
	}

	Routes["create"] = routes.Route{
		Handler: http.HandlerFunc(createUser),
		Methods: []string{"POST"},
	}
}
