package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var buildTargets = map[string]buildTarget{
	"self": {
		description: "Re-build EasyMaker",
		command: func() bool {
			return buildGo("easymaker", "./app/build", ".exe", "", "exe", color.Output, color.Error)
		},
	},
	"cli": {
		description: "Build EasyServer console runner",
		command: func() bool {
			return buildGo("easyserver", "./app/cli", ".exe", "", "exe", color.Output, color.Error)
		},
	},
	"ffi": {
		description: "Build C/C++ bindings",
		command: func() bool {
			return buildGo("easyserver", "./ffi", ".dll", ".so", "c-shared", color.Output, color.Error)
		},
	},
	"config": {
		description: "Build EasyConfig (requires Python3 installed)",
		command: func() bool {
			color.New(color.Bold).Println("- Build C library")
			ok := buildGo("easyserver", "./ffi", ".dll", ".so", "c-shared", color.Output, color.Error)
			if !ok {
				fmt.Fprintln(os.Stderr)
				color.New(color.FgHiRed).Add(color.Bold).Fprintln(color.Error, "Abort")
				return false
			}

			color.New(color.Bold).Println("- Build UI")
			for ui, dst := range pyQtDesign {
				color.New(color.Italic).Println("-- File " + ui + " -> " + dst)
				ok := buildPyQtUI(dst, ui, color.Output, color.Error)
				if !ok {
					fmt.Fprintln(os.Stderr)
					return false
				}
			}

			ext := ".so"
			if runtime.GOOS == "windows" {
				ext = ".dll"
			}
			color.New(color.Bold).Println("- Build executable")
			ok = buildPython("easyconfig", "gui.py", false, map[string]string{"easyserver" + ext: "."}, nil, color.Output, color.Error)
			if !ok {
				fmt.Fprintln(os.Stderr)
			}
			return ok
		},
	},
	"all": {
		description: "Build all targets except \"self\"",
		command: func() bool {
			ok := runSubcommand("CLI", func(stdout, stderr io.Writer) bool {
				return buildGo("easyserver", "./app/cli", ".exe", "", "exe", stdout, stderr)
			})
			if !ok {
				color.New(color.FgHiRed).Add(color.Bold).Fprintln(color.Error, "Abort")
				return false
			}

			ok = runSubcommand("FFI", func(stdout, stderr io.Writer) bool {
				return buildGo("easyserver", "./ffi", ".dll", ".so", "c-shared", stdout, stderr)
			})
			if !ok {
				color.New(color.FgHiRed).Add(color.Bold).Fprintln(color.Error, "Abort")
				return false
			}

			return runSubcommand("EasyConfig", func(stdout, stderr io.Writer) bool {
				color.New(color.Bold).Fprintln(stdout, "- Build UI")
				for ui, dst := range pyQtDesign {
					color.New(color.Italic).Fprintln(stdout, "-- File "+ui+" -> "+dst)
					ok := buildPyQtUI(dst, ui, stdout, stderr)
					if !ok {
						fmt.Fprintln(stderr)
						return false
					}
				}

				ext := ".so"
				if runtime.GOOS == "windows" {
					ext = ".dll"
				}
				color.New(color.Bold).Fprintln(stdout, "- Build executable")
				ok := buildPython("easyconfig", "gui.py", false, map[string]string{"easyserver" + ext: "."}, nil, stdout, stderr)
				if !ok {
					fmt.Fprintln(os.Stderr)
				}
				return ok
			})
		},
	},
}

func runSubcommand(header string, cmd func(stdout, stderr io.Writer) bool) (ok bool) {
	str := color.New(color.Bold).Sprintf("# %s: ", header)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(color.Output))
	s.Prefix = str
	var out bytes.Buffer
	s.Start()
	ok = cmd(&out, &out)
	s.Stop()

	good := color.New(color.FgHiGreen).Sprint("OK")
	bad := color.New(color.FgHiRed).Sprint("Fail")
	if ok {
		fmt.Fprintln(color.Output, str+good)
	} else {
		fmt.Fprintln(color.Output, str+bad)
		io.Copy(color.Error, &out)
	}
	return
}
