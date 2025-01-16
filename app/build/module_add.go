package main

import (
	"errors"
	"flag"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

var moduleAddFlagSet = flag.NewFlagSet("add-module", flag.ExitOnError)
var newModuleLang = moduleAddFlagSet.String("lang", "go", "Language for module (Availible: go, python, c)")
var moduleAddNotInteractive = moduleAddFlagSet.Bool("no-menu", false, "Disable select menu. If not set, all flags are ignored")
var newModuleName = moduleAddFlagSet.String("name", "", "Name of module")
var enabeNewModule = moduleAddFlagSet.Bool("enable", false, "Enable module")

func runModuleAdd() int {
	header := color.New(color.Bold)

	header.Println("Create new module")
	if *moduleAddNotInteractive {
		if createNewModule(*newModuleLang, *newModuleName, *enabeNewModule) {
			return 0
		} else {
			return 1
		}
	}

	existingModules, err := GetModules()
	if err != nil {
		color.New(color.FgHiRed).Add(color.Bold).Fprintf(color.Error, "Error while getting modules: %v\n", err)
	}

	validate := func(input string) error {
		if input == "" {
			return errors.New("empty prompt")
		}

		_, ok := existingModules[input]
		if ok {
			return errors.New("module with this name exists")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Module name",
		Validate: validate,
	}

	moduleName, err := prompt.Run()
	if err != nil {
		color.New(color.FgHiRed).Add(color.Bold).Fprintf(color.Error, "Prompt failed: %v\n", err)
		return -1
	}

	langs := []string{"go", "c", "python"}

	choose := promptui.Select{
		Label: "Choose module language",
		Items: []string{"Go", "C/C++", "Python 3"},
	}

	i, _, err := choose.Run()
	if err != nil {
		color.New(color.FgHiRed).Add(color.Bold).Fprintf(color.Error, "Prompt failed: %v\n", err)
		return -1
	}
	lang := langs[i]

	choose = promptui.Select{
		Label: "Do you want to enable new module?",
		Items: []string{"Yes", "No"},
	}

	i, _, err = choose.Run()
	if err != nil {
		color.New(color.FgHiRed).Add(color.Bold).Fprintf(color.Error, "Prompt failed: %v\n", err)
		return -1
	}
	enable := i == 0

	if createNewModule(lang, moduleName, enable) {
		return 0
	} else {
		return 1
	}
}

func createNewModule(lang string, name string, enable bool) bool {
	errColor := color.New(color.FgHiRed).Add(color.Bold)
	home, err := filepath.Abs(".")
	if err != nil {
		errColor.Fprintf(color.Error, "Cannot get absolute path current dictory: %v\n", err)
		return false
	}

	err = os.Mkdir(filepath.Join(home, "modules", name), 0o755)
	if err != nil {
		errColor.Fprintf(color.Error, "Cannot create module dir: %v\n", err)
		return false
	}

	tmpl, err := template.ParseFiles(filepath.Join(home, "templates", "build", lang+".export.txt"))
	if err != nil {
		errColor.Fprintf(color.Error, "Error when parsing export template: %v\n", err)
		return false
	}

	export, err := os.Create(filepath.Join(home, "modules", name, "export.go"))
	if err != nil {
		errColor.Fprintf(color.Error, "Cannot create `export.go` file: %v\n", err)
		return false
	}
	defer export.Close()

	err = tmpl.Execute(export, name)
	if err != nil {
		errColor.Fprintf(color.Error, "Error when setting export file: %v\n", err)
		return false
	}

	if enable {
		return enableModule(name)
	}

	return true
}

func enableModule(name string) bool {
	errColor := color.New(color.FgHiRed).Add(color.Bold)
	modules, err := GetModules()
	if err != nil {
		errColor.Fprintf(color.Error, "Error when getting enabled modules: %v\n", err)
		return false
	}

	modules[name] = true
	active := make([]string, 0, len(modules))

	i := 0
	for module, isActive := range modules {
		if isActive {
			active = append(active, module)
			i++
		}
	}

	err = WriteActiveModules(active)
	if err != nil {
		errColor.Fprintf(color.Error, "Error when saving enabled modules: %v\n", err)
		return false
	}

	return true
}

func GetModules() (map[string]bool, error) {
	res := make(map[string]bool)

	home, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	modules, err := filepath.Glob(filepath.Join(home, "modules", "*"))
	if err != nil {
		return nil, err
	}

	for _, module := range modules {
		name := filepath.Base(module)
		res[name] = false
	}

	fileSet := token.NewFileSet()

	f, err := parser.ParseFile(fileSet, filepath.Join(home, "config", "modules", "imports.go"), nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	for _, imp := range f.Imports {
		name := path.Base(imp.Path.Value)
		res[name[:len(name)-1]] = true
	}
	return res, nil
}

func WriteActiveModules(modules []string) error {
	home, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(home, "config", "modules", "imports.go"), os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("package modules\n\n")
	if err != nil {
		return err
	}
	if len(modules) == 0 {
		return nil
	}

	if len(modules) == 1 {
		_, err = file.WriteString("import _ \"github.com/BIQ-Cat/easyserver/modules/" + modules[0] + "\"\n")
		if err != nil {
			return err
		}
		return nil
	}

	_, err = file.WriteString("import (\n")
	if err != nil {
		return err
	}

	for _, module := range modules {
		_, err = file.WriteString("\t_ \"github.com/BIQ-Cat/easyserver/modules/" + module + "\"\n")
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(")\n")

	return err
}
