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
	t.OperationTypes = append(types.OperationTypes, types.InvalidOpType)
	// Optimism does not have miner or uncle reward type ops
	// TODO(inphi): add BurnOpType
	t.OperationTypes = t.OperationTypes[2:]
	t.OperationTypes = append(types.OperationTypes, types.InvalidOpType, MintOpType)

	client, err := NewOpClient(cfg)
	if err != nil {
		log.Fatalf("cannot initialize client: %v", err)
	}

	err = utils.BootStrap(cfg, t, types.Errors, client)
	if err != nil {
		log.Fatalf("unable to bootstrap Rosetta server: %v", err)
	}
}
