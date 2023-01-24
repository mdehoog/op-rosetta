package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
)

// MintOps constructs a list of [RosettaTypes.Operation]s for an Optimism Deposit or "mint" transaction.
func MintOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	if tx.Transaction.Mint() == nil {
		return nil
	}
	return []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(startIndex),
			},
			Type:   MintOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: tx.From.String(),
			},
			Amount: evmClient.Amount(tx.Transaction.Mint(), sdkTypes.Currency),
		},
	}
}
