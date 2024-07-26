package auth

import (
	"github.com/BIQ-Cat/easyserver/routes"
)

var Routes map[string]routes.Route

func init() {
	routes.Routes["auth"] = &Routes
}
