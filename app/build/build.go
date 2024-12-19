package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

type buildTarget struct {
	description string
	command     func() bool
}

var buildFlagSet = flag.NewFlagSet("build", flag.ExitOnError)
var buildCurrentTarget = buildFlagSet.String("target", "all", "Target to build. Use special target 'help' to get information about availibe targets.")
var buildOutput = buildFlagSet.String("out", "./build", "Build destination")

func runBuild() {
	header := color.New(color.Bold)
	header.Println("Build")
	if *buildCurrentTarget == "help" {
		header.Println("Availible targets")
		for name, target := range buildTargets {
			fmt.Fprintf(color.Output, "  %s %s\n", header.Sprint(name+":"), color.New(color.Italic).Add(color.FgBlue).Sprint(target.description))
		}
		return
	}
	for name, target := range buildTargets {
		if *buildCurrentTarget == name {
			target.command()
			color.New(color.FgHiGreen).Add(color.Bold).Println("Success!")
			return
		}
	}
	color.New(color.FgHiRed).Add(color.Bold).Fprintln(colorable.NewColorable(os.Stderr), "Undefined target")
}
