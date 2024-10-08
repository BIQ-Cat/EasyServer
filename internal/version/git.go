package main

import "fmt"

var gitTag = "0.01dev"
var gitCommit = "b86d52"
var gitBranch = "main"

func GetNightVersion() string {
	return fmt.Sprintf("%s (%s)", gitCommit, gitBranch)
}

func GetVersion() string {
	return gitTag
}
