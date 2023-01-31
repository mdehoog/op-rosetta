package utils

import (
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	common "github.com/mdehoog/op-rosetta/pkg/common"
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
	ots = append(ots, common.MintOpType)
	// ots = append(ots, handlers.BurnOpType)
	ots = append(ots, common.StopOpType)
	t.OperationTypes = ots

	return t
}
