package main

import (
	"log"

	"github.com/coinbase/rosetta-geth-sdk/types"
	"github.com/coinbase/rosetta-geth-sdk/utils"
)

func main() {
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
	ots = append(ots, types.InvalidOpType)
	ots = append(ots, MintOpType)
	// TODO(inphi): add BurnOpType
	t.OperationTypes = ots

	client, err := NewOpClient(cfg)
	if err != nil {
		log.Fatalf("cannot initialize client: %v", err)
	}

	err = utils.BootStrap(cfg, t, types.Errors, client)
	if err != nil {
		log.Fatalf("unable to bootstrap Rosetta server: %v", err)
	}
}
