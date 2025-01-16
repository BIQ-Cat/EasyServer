package main

import "flag"

type command struct {
	run         func() int
	description string
	flagSet     *flag.FlagSet
}

var commands = map[string]command{
	"build": {
		run:         runBuild,
		description: "builds specified target",
		flagSet:     buildFlagSet,
	},

	"add-module": {
		run:         runModuleAdd,
		description: "Adds a new module",
		flagSet:     moduleAddFlagSet,
	},
}
