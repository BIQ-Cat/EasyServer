package utils

import (
	"fmt"
	"log"
)

var debug bool

func SetDebug(state bool) {
	debug = state
}

func LogInDebug(prompt ...any) {
	if debug {
		fmt.Println(prompt...)
	}
}

func LogInProduction(prompt ...any) {
	if debug {
		log.Println(prompt...)
	}
}

func CheckDebugState() bool { return debug }
