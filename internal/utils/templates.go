package utils

import (
	"html/template"
	"os"
	"path/filepath"
)

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string

	templateDir := filepath.Join("templates", dir)
	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
