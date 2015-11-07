package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	"github.com/op/go-logging"
)

var (
	commit  string
	version string
	build   string
	logger  *logging.Logger
)

const logFormat = "%{color}%{level}%{color:reset} %{message}"

func init() {
	logging.SetFormatter(logging.MustStringFormatter(logFormat))
	logger = logging.MustGetLogger("fetcher")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	parser := flags.NewParser(nil, flags.Default)
	parser.AddCommand("proxy", "", "", &ProxyCommand{})
	parser.AddCommand("server", "", "", &ServerCommand{})

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			parser.WriteHelp(os.Stdout)
			fmt.Printf("\nBuild information\n  version: %s\n  build: %s\n  commit: %s\n", version, build, commit)
		}

		os.Exit(1)
	}
}
