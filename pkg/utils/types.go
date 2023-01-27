package utils

import (
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"

	"github.com/mdehoog/op-rosetta/pkg/handlers"
)

// LoadTypes loads the types for the Optimism Rosetta service.
func LoadTypes() *sdkTypes.Types {
	t := sdkTypes.LoadTypes()
	var ots []string
	for _, ot := range t.OperationTypes {
		if ot == sdkTypes.MinerRewardOpType || ot == sdkTypes.UncleRewardOpType {
			// Optimism does not have miner or uncle reward type ops
			continue
		}
		ots = append(ots, ot)
	}
	ots = append(ots, handlers.MintOpType)
	t.OperationTypes = ots

	return t
}
