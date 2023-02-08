package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	common "github.com/mdehoog/op-rosetta/pkg/common"
)

// MintOps constructs a list of [RosettaTypes.Operation]s for an Optimism Deposit or "mint" transaction.
func MintOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	if tx.Transaction.Mint() == nil {
		return nil
	}

	opIndex := int64(startIndex)
	opType := common.MintOpType
	opStatus := sdkTypes.SuccessStatus
	toAddress := evmClient.MustChecksum(tx.Transaction.To().String())
	amount := evmClient.Amount(tx.Transaction.Mint(), sdkTypes.Currency)

	return []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: opIndex,
			},
			Type:   opType,
			Status: RosettaTypes.String(opStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: toAddress,
			},
			Amount: amount,
		},
	}
}
