package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	common "github.com/mdehoog/op-rosetta/pkg/common"
)

// BurnOps constructs a list of [RosettaTypes.Operation]s for an Optimism Withdrawal or "burn" transaction.
func BurnOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	// tx's To address could be nil, if the type == CREATE, it will break
	if *tx.Transaction.To() != common.L2ToL1MessagePasser {
		return nil
	}

	opIndex := int64(startIndex)
	opType := common.BurnOpType
	opStatus := sdkTypes.SuccessStatus
	fromAddress := evmClient.MustChecksum(tx.From.String())
	amount := evmClient.Amount(tx.Transaction.Value(), sdkTypes.Currency)

	return []*RosettaTypes.Operation{
		GenerateOp(opIndex, nil, opType, opStatus, fromAddress, amount, nil),
	}
}
