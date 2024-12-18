package easyserver

import (
	"fmt"
	"log"
	"net/http"
	pathlib "path"
	"slices"
	"strings"

	config "github.com/BIQ-Cat/easyserver/config/base"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Enable postgres driver
)

const (
	RouteNotFoundPath = "?404"
)

type Router struct {
	Modules    map[string]Module
	Default404 http.Handler
	db         *gorm.DB
}

func (router *Router) sortedModulePaths() []string {
	list := make([]string, len(router.Modules))
	i := 0
	for key := range router.Modules {
		list[i] = key
		i++
	}
	slices.SortFunc(list, func(a string, b string) int {
		return len(strings.Split(strings.TrimPrefix(b, "/"), "/")) - len(strings.Split(strings.TrimPrefix(a, "/"), "/"))
	})
	return list
}

func (router *Router) GetController(requestPath string) *Controller {
	for path, modules := range router.Modules {
		for subpath, controller := range modules.Route {
			fullPath := fmt.Sprintf("/%s/%s", path, subpath)
			if fullPath == requestPath {
				return &controller
			}
		}
	}

	if config.Config.Debug {
		log.Println("WARNING: no controller found")
	}
	return nil // 404
}

type Module struct {
	Route
	Middlewares []MiddlewareFunc
	NotExplict  bool
	Models      []interface{}
}

func (module Module) processController(path string, w http.ResponseWriter, r *http.Request) (ok bool) {
	for _, controllerPath := range module.Route.sortedControllerPaths() {
		controller := module.Route[controllerPath]
		if !isControllerEnabled(r, controller) {
			continue
		}
		if module.NotExplict || strings.HasSuffix(controllerPath, "*") {
			notExplictPath := strings.TrimSuffix(controllerPath, "*")
			fullPath := pathlib.Join(path, notExplictPath)
			if strings.HasPrefix(r.URL.Path, fullPath) {
				controller.ServeHTTP(w, r)
				return true
			}
		} else {
			if pathlib.Join(path, controllerPath) == r.URL.Path {
				controller.ServeHTTP(w, r)
				return true
			}
		}
	}
	return false
}

type Route map[string]Controller

func (route Route) UpdateRouteNotFoundController(controller Controller) {
	route[RouteNotFoundPath] = controller
}

func (route Route) sortedControllerPaths() []string {
	list := make([]string, len(route))
	i := 0
	for key := range route {
		list[i] = key
		i++
	}
	slices.SortFunc(list, func(a string, b string) int {
		return len(strings.Split(strings.TrimPrefix(b, "/"), "/")) - len(strings.Split(strings.TrimPrefix(a, "/"), "/"))
	})
	return list
}

type Controller struct {
	http.Handler
	Methods []string
	Headers map[string]string
	Schemas []string
	Data    map[string]any
}

type MiddlewareFunc func(next http.Handler) http.Handler

func (router *Router) NoMiddlewares(w http.ResponseWriter, r *http.Request) {
	newRouter := router.withoutMiddlewares()
	newRouter.ServeHTTP(w, r)
}

func (router *Router) withoutMiddlewares() http.Handler {
	newRouter := new(Router)
	newRouter.Modules = make(map[string]Module)
	for path, module := range router.Modules {
		newRouter.Modules[path] = Module{
			Route:       module.Route,
			Middlewares: []MiddlewareFunc{},
			Models:      module.Models,
		}
	}

	return newRouter
}

func (router *Router) Connect() error {
	username := config.EnvConfig.DBUser
	password := config.EnvConfig.DBPass
	dbName := config.EnvConfig.DBName
	dbHost := config.EnvConfig.DBHost
	dbPort := config.EnvConfig.DBPort

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=disable password=%s", dbHost, username, dbName, dbPort, password) //Создать строку подключения

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		return err
	}

	router.db = conn
	if config.Config.Debug {
		router.db.Debug()
		for _, module := range router.Modules {
			router.db.AutoMigrate(module.Models...)
		}
	}
	return nil
}

func (router *Router) DB() *gorm.DB {
	return router.db
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	var gotMiddleware bool
	handler = router.withoutMiddlewares()
	for path, module := range router.Modules {
		if strings.HasPrefix(r.URL.Path, path) {
			for _, middlew := range module.Middlewares {
				gotMiddleware = true
				handler = middlew(handler)
			}
		}
	}
	if gotMiddleware {
		handler.ServeHTTP(w, r)
		return
	}

	page404 := router.Default404
	got404Earlier := false
	if page404 == nil {
		page404 = http.NotFoundHandler()
	}
	for _, path := range router.sortedModulePaths() {
		module := router.Modules[path]
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		if strings.HasPrefix(r.URL.Path, path) {
			if module.processController(path, w, r) {
				return
			} else {
				notFoundHandler, gotNotFoundHandler := module.Route[RouteNotFoundPath]
				if gotNotFoundHandler && !got404Earlier {
					got404Earlier = true
					page404 = notFoundHandler
				}
			}
		}
	}
	page404.ServeHTTP(w, r)
}

func isControllerEnabled(r *http.Request, controller Controller) bool {
	for header, value := range controller.Headers {
		if r.Header.Get(header) != value {
			return false
		}
	}
	if len(controller.Methods) != 0 && !slices.Contains(controller.Methods, r.Method) || len(controller.Schemas) != 0 && !slices.Contains(controller.Schemas, r.URL.Scheme) {
		return false
	}

	return true
}
