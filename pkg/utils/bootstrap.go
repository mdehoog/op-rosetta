package utils

import (
	"log"

	"github.com/mdehoog/op-rosetta/pkg/client"
	"github.com/mdehoog/op-rosetta/pkg/config"

	rosettaGethTypes "github.com/coinbase/rosetta-geth-sdk/types"
	rosettaGethUtils "github.com/coinbase/rosetta-geth-sdk/utils"
	"github.com/urfave/cli"
)

// Bootstrap returns a closure that executes the core logic of op-rosetta.
// The closure constructs a new OpClient and runs the Rosetta server using [rosettaGethUtils.Bootstrap].
// The use of a closure allows the parameters bound to the top-level main package, e.g.
// GitVersion, to be captured and used once the function is executed.
func Bootstrap() func(cliCtx *cli.Context) error {
	return func(cliCtx *cli.Context) error {
		// Load the [configuration.Configuration]
		cfg, err := config.LoadConfiguration()
		if err != nil {
			log.Fatalf("unable to load configuration: %v", err)
		}

		// Load [types.Types]
		t := LoadTypes()

		// Construct a new OpClient with handlers
		client, err := client.NewOpClient(cfg)
		if err != nil {
			log.Fatalf("cannot initialize client: %v", err)
		}

		// Boostrap the Rosetta server, serving RESTful requests
		// This method is a blocking call provided by the rosetta-geth-sdk
		err = rosettaGethUtils.BootStrap(cfg, t, rosettaGethTypes.Errors, client)
		if err != nil {
			log.Fatalf("unable to bootstrap Rosetta server: %v", err)
		}

		return nil
	}
}
