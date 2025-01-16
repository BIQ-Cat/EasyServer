package python

import (
	"path/filepath"

	"github.com/BIQ-Cat/easyserver/internal/api"
	"github.com/BIQ-Cat/easyserver/internal/router"
)

func init() {
	api.StartPython()
	defer api.EndPython()

	path, err := filepath.Abs(filepath.Join(".", "modules", "python"))
	if err != nil {
		panic(err)
	}

	api.PythonImportLib()

	module, ok := api.CreateModule(path, "python")
	if !ok {
		panic("bad Python")
	}

	router.DefaultRouter.Modules["python"] = module
}
