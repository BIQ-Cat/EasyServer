package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/BIQ-Cat/easyserver"
	"github.com/BIQ-Cat/easyserver/internal/api"
)

func main() {
	api.StartPython()
	defer api.EndPython()

	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	pLib := api.PythonImport(path, "lib")
	pModule := api.PythonImport(filepath.Join(path, "modules", "python"), "main")
	module, ok := api.CreateModule(pModule, pLib)
	if !ok {
		return
	}
	fmt.Println(module)

	router := easyserver.Router{
		Modules: map[string]easyserver.Module{
			"python": module,
		},
	}

	http.ListenAndServe(":8080", &router)
}
