package routes

import (
	"fmt"
	"log"

	config "github.com/BIQ-Cat/easyserver/config/base"
)

func GetController(requestPath string) *Controller {
	for path, subroutes := range Routes {
		for subpath, route := range *subroutes {
			fullPath := fmt.Sprintf("/%s/%s", path, subpath)
			if fullPath == requestPath {
				return &route
			}
		}
	}

	if config.Config.Debug {
		log.Println("WARNING: no controller found")
	}
	return nil // 404
}
