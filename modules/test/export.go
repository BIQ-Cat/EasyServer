package test

import (
	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/internal/router"
)

func init() {
	router.DefaultRouter.Modules["test"] = easyserver.Module{

	}

	// var cfg basictypes.JSONConfig = moduleconfig.Config
	// json.Configurations["test"] = &cfg
}
