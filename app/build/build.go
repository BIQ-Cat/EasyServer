package main

import (
	"flag"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/fatih/color"
)

type buildTarget struct {
	description string
	command     func() bool
}

var buildFlagSet = flag.NewFlagSet("build", flag.ExitOnError)
var buildCurrentTarget = buildFlagSet.String("target", "all", "Target to build. Use special target 'help' to get information about availibe targets.")
var buildOutput = buildFlagSet.String("out", "./build", "Build destination")

func runBuild() int {
	header := color.New(color.Bold)
	header.Println("Build")
	if *buildCurrentTarget == "help" {
		header.Println("Availible targets")
		for name, target := range buildTargets {
			fmt.Fprintf(color.Output, "  %s %s\n", header.Sprint(name+":"), color.New(color.Italic).Add(color.FgBlue).Sprint(target.description))
		}
		return 0
	}
	for name, target := range buildTargets {
		if *buildCurrentTarget == name {
			if target.command() {
				color.New(color.FgHiGreen).Add(color.Bold).Println("Success!")
				return 0
			} else {
				color.New(color.FgHiRed).Add(color.Bold).Println("Failure!")
				return 1
			}
		}
	}
	color.New(color.FgHiRed).Add(color.Bold).Fprintln(color.Error, "Undefined target")
	return -1
}

func buildGo(dst, src string, windowsExt, unixExt string, buildmode string, stdout, stderr io.Writer) (ok bool) {
	ext := unixExt
	if runtime.GOOS == "windows" {
		ext = windowsExt
	}
	errorColor := color.New(color.FgHiRed).Add(color.Bold)
	path, err := exec.LookPath("go")
	if err != nil {
		errorColor.Fprintln(stderr, "Go executable not found")
		return false
	}

	cmd := exec.Command(path, "build", "-buildmode="+buildmode, "-linkshared", "-o", filepath.Join(*buildOutput, dst+ext), src)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return !processCmdRunError(cmd.Run(), stderr)
}

func buildPyQtUI(dst string, ui string, stdout, stderr io.Writer) (ok bool) {
	pyuic6, err := exec.LookPath("pyuic6")
	errorColor := color.New(color.FgHiRed).Add(color.Bold)
	if err != nil {
		errorColor.Fprintln(stderr, "`pyuic` wasn't found in PATH. Consider installing PyQt6 or entering python venv.")
		return false
	}

	dir := filepath.Join(".", "app", "python", "ui")
	cmd := exec.Command(pyuic6, "-o", filepath.Join(dir, dst), filepath.Join(dir, ui))
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	return !processCmdRunError(cmd.Run(), stderr)
}

func buildPython(name string, src string, console bool, binaries map[string]string, data map[string]string, stdout, stderr io.Writer) (ok bool) {
	pyinstaller, err := exec.LookPath("pyinstaller")
	errorColor := color.New(color.FgHiRed).Add(color.Bold)

	if err != nil {
		errorColor.Fprintln(stderr, "PyInstaller wasn't found in PATH. Consider installing it or entering python venv.")
		return false
	}

	consoleFlag := "--noconsole"
	if console {
		consoleFlag = "--console"
	}

	binaryFlags := make([]string, len(binaries)*2)
	i := 0

	for src, dst := range binaries {
		binaryFlags[i] = "--add-binary"
		binaryFlags[i+1] = filepath.Join("..", src) + ":" + dst
		i += 2
	}

	dataFlags := make([]string, len(data)*2)
	i = 0

	for src, dst := range data {
		dataFlags[i] = "--add-data"
		dataFlags[i+1] = filepath.Join(".", "app", "python", src) + ":" + dst
		i += 2
	}

	args := append(
		slices.Concat(
			[]string{
				"--distpath", *buildOutput,
				"--workpath", filepath.Join(*buildOutput, "pyinstaller"),
				"--specpath", filepath.Join(*buildOutput, "pyinstaller"),
				"--name", name,
				"--paths", "app/python",
				"--onefile",
				consoleFlag,
				"--clean",
				"--noconfirm",
				"--log-level", "WARN",
			},
			dataFlags, binaryFlags),
		filepath.Join(".", "app", "python", src))

	cmd := exec.Command(pyinstaller, args...)

	cmd.Stderr = stderr
	cmd.Stdout = stdout

	return !processCmdRunError(cmd.Run(), stderr)
}

func processCmdRunError(err error, stderr io.Writer) (wasError bool) {
	errorColor := color.New(color.FgHiRed).Add(color.Bold)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if !exitErr.Success() {
				errorColor.Fprintf(stderr, "Got error while building (code %d). Abort", exitErr.ExitCode())
				return true
			}
		} else {
			errorColor.Fprintln(stderr, "Command execution error!")
			color.New(color.Bold).Fprintln(stderr, err)
			return true
		}
	}
	return false
}
