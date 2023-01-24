// Package app contains the primary functionality of op-rosetta.
package app

import (
	"log"

	"github.com/coinbase/rosetta-geth-sdk/types"
	"github.com/coinbase/rosetta-geth-sdk/utils"
	"github.com/mdehoog/op-rosetta/handlers"
	"github.com/urfave/cli"
)

// Main is the entrypoint into op-rosetta. This method returns a
// closure that executes the service and blocks until the service exits. The use
// of a closure allows the parameters bound to the top-level main package, e.g.
// GitVersion, to be captured and used once the function is executed.
func Main(version string) func(cliCtx *cli.Context) error {
	return func(cliCtx *cli.Context) error {
		cfg, err := LoadConfiguration()
		if err != nil {
			log.Fatalf("unable to load configuration: %v", err)
		}

		t := types.LoadTypes()
		var ots []string
		for _, ot := range t.OperationTypes {
			if ot == types.MinerRewardOpType || ot == types.UncleRewardOpType {
				// Optimism does not have miner or uncle reward type ops
				continue
			}
			ots = append(ots, ot)
		}
		ots = append(ots, handlers.MintOpType)
		ots = append(ots, handlers.BurnOpType)
		t.OperationTypes = ots

		client, err := NewOpClient(cfg)
		if err != nil {
			log.Fatalf("cannot initialize client: %v", err)
		}

		err = utils.BootStrap(cfg, t, types.Errors, client)
		if err != nil {
			log.Fatalf("unable to bootstrap Rosetta server: %v", err)
		}

		return nil
	}
}
