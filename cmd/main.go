package main

import (
	"fmt"
	"os"

	"github.com/mdehoog/op-rosetta/app"
	"github.com/mdehoog/op-rosetta/utils"
	"github.com/urfave/cli"

	"github.com/ethereum/go-ethereum/log"
)

var (
	Version   = "v0.1.0"
	GitCommit = ""
	GitDate   = ""
)

func main() {
	utils.SetupDefaults()

	cliApp := cli.NewApp()
	cliApp.Flags = []cli.Flag{}
	cliApp.Version = fmt.Sprintf("%s-%s-%s", Version, GitCommit, GitDate)
	cliApp.Name = "op-rosetta"
	cliApp.Usage = "Optimism Rosetta Service"
	cliApp.Description = "Service for translating Optimism transactions into Rosetta Operations"

	cliApp.Action = app.Main(Version)
	err := cliApp.Run(os.Args)
	if err != nil {
		log.Crit("op-rosetta failed", "message", err)
	}
}
