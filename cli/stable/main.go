package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
)

var (
	commit  string
	version string
	build   string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	parser := flags.NewParser(nil, flags.Default)
	parser.AddCommand("server", "", "", &ServerCommand{})

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			parser.WriteHelp(os.Stdout)
			fmt.Printf(
				"\nBuild information\n  version: %s\n  build: %s\n  commit: %s\n",
				version, build, commit,
			)
		}

		fmt.Fprintf(os.Stdout, err.Error())
		os.Exit(1)
	}
}
