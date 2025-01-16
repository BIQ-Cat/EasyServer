package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var noColor = flag.Bool("no-color", false, "Disable color output")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: easymaker <flags> [COMMAND] ...")
		fmt.Fprintln(os.Stderr, "Flags: ")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Commands: ")
		for name, command := range commands {
			fmt.Fprintln(os.Stderr, name+": "+command.description)
			command.flagSet.PrintDefaults()
		}
	}
	flag.Parse()

	if *noColor {
		color.NoColor = true
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	for name, cmd := range commands {
		if args[0] == name {
			cmd.flagSet.Parse(args[1:])
			exitCode := cmd.run()
			os.Exit(exitCode)
		}
	}

	fmt.Fprintln(os.Stderr, "Command \""+args[0]+"\" not found")
	flag.Usage()
	os.Exit(-1)
}
