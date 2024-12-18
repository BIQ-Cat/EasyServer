package router

import "github.com/BIQ-Cat/easyserver"

var DefaultRouter = easyserver.Router{
	Modules: make(map[string]easyserver.Module),
}
