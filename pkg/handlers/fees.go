// Package handlers handles translating transactions into [RosettaTypes.Operation]s.
package handlers

import (
	"math/big"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	common "github.com/mdehoog/op-rosetta/pkg/common"
)

// FeeOps returns the fee operations for a given transaction.
func FeeOps(tx *evmClient.LoadedTransaction) ([]*RosettaTypes.Operation, error) {
	if tx.Transaction.IsDepositTx() {
		return nil, nil
	}

	var receipt EthTypes.Receipt
	if err := receipt.UnmarshalJSON(tx.Receipt.RawMessage); err != nil {
		return nil, err
	}

	sequencerFeeAmount := new(big.Int).Set(tx.FeeAmount)
	if tx.FeeBurned != nil {
		sequencerFeeAmount.Sub(sequencerFeeAmount, tx.FeeBurned)
	}
	if receipt.L1Fee != nil {
		sequencerFeeAmount.Sub(sequencerFeeAmount, receipt.L1Fee)
	}

	if sequencerFeeAmount == nil {
		return nil, nil
	}

	feeRewarder := tx.Miner
	if len(tx.Author) > 0 {
		feeRewarder = tx.Author
	}

	ops := []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: 0,
			},
			Type:   sdkTypes.FeeOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: evmClient.MustChecksum(tx.From.String()),
			},
			Amount: evmClient.Amount(new(big.Int).Neg(tx.Receipt.TransactionFee), sdkTypes.Currency),
		},

		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: 1,
			},
			RelatedOperations: []*RosettaTypes.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Type:   sdkTypes.FeeOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: evmClient.MustChecksum(feeRewarder),
			},
			Amount: evmClient.Amount(sequencerFeeAmount, sdkTypes.Currency),
		},

		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: 2,
			},
			RelatedOperations: []*RosettaTypes.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Type:   sdkTypes.FeeOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: common.BaseFeeVault.Hex(),
			},
			// Note: The basefee is not actually burned on L2
			Amount: evmClient.Amount(tx.FeeBurned, sdkTypes.Currency),
		},

		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: 3,
			},
			RelatedOperations: []*RosettaTypes.OperationIdentifier{
				{
					Index: 0,
				},
			},
			Type:   sdkTypes.FeeOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: common.L1FeeVault.Hex(),
			},
			Amount: evmClient.Amount(receipt.L1Fee, sdkTypes.Currency),
		},
	}

	return ops, nil
}
