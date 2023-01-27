package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
)

// BurnOps constructs a list of [RosettaTypes.Operation]s for an Optimism Withdrawal or "burn" transaction.
func BurnOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	if *tx.Transaction.To() != L2ToL1MessagePasser {
		return nil
	}
	return []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(startIndex),
			},
			Type:   BurnOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: tx.From.String(),
			},
			Amount: evmClient.Amount(tx.Transaction.Value(), sdkTypes.Currency),
		},
	}
}
