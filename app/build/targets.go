package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var buildTargets = map[string]buildTarget{
	"self": {
		description: "Re-build EasyMaker",
		command: func() bool {
			return buildGo("easymaker", "./app/build", ".exe", "")
		},
	},
	"cli": {
		description: "Build EasyServer console runner",
		command: func() bool {
			return buildGo("easyserver", "./app/cli", ".exe", "")
		},
	},
	"ffi": {
		description: "Build FFI bindings",
		command: func() bool {
			return buildGo("easyserver", "./ffi", ".dll", ".so")
		},
	},
	"config": {
		description: "Build EasyConfig (requires Python3 installed)",
		command: func() bool {
			return true
		},
	},
}

func buildGo(dst, src string, windowsExt, unixExt string) (ok bool) {
	ext := unixExt
	if runtime.GOOS == "windows" {
		ext = windowsExt
	}
	errorColor := color.New(color.FgHiRed).Add(color.Bold)
	stderr := colorable.NewColorable(os.Stderr)
	path, err := exec.LookPath("go")
	if err != nil {
		errorColor.Fprintln(stderr, "Go executable not found")
		return false
	}
	cmd := exec.Command(path, "build", "-o", filepath.Join(*buildOutput, dst+ext), src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if !exitErr.Success() {
				errorColor.Fprintf(stderr, "Got error while building (code %d). Abort", exitErr.ExitCode())
				return false
			}
		} else {
			errorColor.Fprintln(stderr, "Command execution error!")
			color.New(color.Bold).Fprintln(stderr, err)
			return false
		}
	}
	return true
}
