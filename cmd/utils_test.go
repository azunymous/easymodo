package cmd

import "os"

var wd string

func init() {
	wd, _ = os.Getwd()
}

func cleanup() {
	_ = os.Chdir(wd)
}
