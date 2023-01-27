package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/log"
)

// MintOps constructs a list of [RosettaTypes.Operation]s for an Optimism Deposit or "mint" transaction.
func MintOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	if tx.Transaction.Mint() == nil {
		return nil
	}
	log.Info("mint operation detected", "tx", tx.TxHash)
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
