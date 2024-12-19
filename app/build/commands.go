package main

import "flag"

type command struct {
	run         func()
	description string
	flagSet     *flag.FlagSet
}

var commands = map[string]command{
	"build": {
		run:         runBuild,
		description: "builds specified target",
		flagSet:     buildFlagSet,
	},
}
