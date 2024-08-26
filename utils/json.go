package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BIQ-Cat/easyserver/config"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func GetStatus(resp map[string]interface{}) bool {
	return resp["status"].(bool)
}

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	if config.Config.Debug {
		fmt.Fprint(w, err)
	}
}
