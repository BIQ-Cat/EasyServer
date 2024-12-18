package moduleconfig

import (
	"net/http"

	"github.com/BIQ-Cat/easyserver/config/modules/cors/types"
)

var Config = types.Config{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodHead,
	},
	AllowedHeaders:   []string{"*"},
	AllowCredentails: false,
}
