package main

import (
	"fmt"
	"gvm"
	"os"

	"gvm/util"
)

func main() {
	parseArgs()
}

func usage() {
	util.PrintlnExit(`
Current gvm version: %s

Usage: 
	gvm <command> <argument>

Example:
	install
		gvm install X.Y.Z
	
	use
		gvm use X.Y.Z
	
	download
		gvm download  X.Y.Z
	
	list
		gvm list
	
	mirror
		gvm mirror
		gvm mirror url
	
	setup
		gvm setup
`, gvm.VERSION)
}

func parseArgs() {
	var args = os.Args[1:]
	if len(args) == 0 {
		usage()
	}

	var err error
	action := args[0]
	switch action {
	case "install", "use", "download":
		if len(args) != 2 {
			usage()
		}
		v := args[1]
		if !gvm.VerifyVersion(v) {
			util.PrintlnExit("invalid version")
		}

		if action == "install" {
			err = gvm.Install(v)
		} else if action == "use" {
			err = gvm.Use(v)
		} else {
			_, err = gvm.Download(v)
		}
	case "list":
		err = gvm.List()
	case "mirror":
		if len(args) == 1 {
			gvm.Mirror()
		} else {
			err = gvm.SetMirror(args[1])
		}
	case "setup":
		err = gvm.Setup()
	default:
		fmt.Println("invalid command")
		usage()
	}

	if err != nil {
		util.PrintlnExit(err.Error())
	}
}
