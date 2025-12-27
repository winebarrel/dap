package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/winebarrel/dap"
)

var (
	version string
)

func parseArgs() *dap.Options {
	var cli struct {
		dap.Options
		Version kong.VersionFlag
	}

	parser := kong.Must(&cli, kong.Vars{"version": version})
	parser.Model.HelpFlag.Help = "Show help."
	_, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	return &cli.Options
}

func main() {
	options := parseArgs()
	server := dap.NewServer(options)
	err := server.Run()

	if err != nil {
		log.Fatal(err)
	}
}
