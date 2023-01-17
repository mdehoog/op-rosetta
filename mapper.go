package main

import (
	"math/big"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"

	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
)

var (
	baseFeeVault = common.HexToAddress("0x4200000000000000000000000000000000000019")
	l1FeeVault   = common.HexToAddress("0x420000000000000000000000000000000000001a")

	MintOpType = "MINT"
)

// FeeOps returns the fee operations for a given transaction
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
				Address: baseFeeVault.Hex(),
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
				Address: l1FeeVault.Hex(),
			},
			Amount: evmClient.Amount(receipt.L1Fee, sdkTypes.Currency),
		},
	}

	return ops, nil
}

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
