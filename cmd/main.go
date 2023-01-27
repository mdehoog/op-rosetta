package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdehoog/op-rosetta/pkg/logging"
	"github.com/mdehoog/op-rosetta/pkg/utils"

	"github.com/urfave/cli"
)

var (
	Version   = "v1.0.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	logging.SetupLogging()

	// Setup App Metadata
	cliApp := cli.NewApp()
	cliApp.Flags = []cli.Flag{}
	cliApp.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	cliApp.Name = "op-rosetta"
	cliApp.Usage = "Optimism Rosetta Service"
	cliApp.Description = "Service for translating Optimism transactions into Rosetta Operations"

	// The main action of the app
	cliApp.Action = utils.Bootstrap()

	// Run op-rosetta
	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatalf("op-rosetta failed: %v", err)
	}
}
