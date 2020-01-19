package cmd

import "os"

var wd string

func init() {
	wd, _ = os.Getwd()
}

func cleanup() {
	w = os.Stdout
	_ = os.Chdir(wd)
}
