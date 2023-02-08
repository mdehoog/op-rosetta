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

	opType := sdkTypes.FeeOpType
	opStatus := sdkTypes.SuccessStatus
	fromAddress := evmClient.MustChecksum(tx.From.String())
	fromAmount := evmClient.Amount(new(big.Int).Neg(tx.Receipt.TransactionFee), sdkTypes.Currency)
	sequencerRelatedOps := []*RosettaTypes.OperationIdentifier{
		{
			Index: 0,
		},
	}
	sequencerAddress := evmClient.MustChecksum(feeRewarder)
	sequencerAmount := evmClient.Amount(sequencerFeeAmount, sdkTypes.Currency)
	baseFeeVaultRelatedOps := []*RosettaTypes.OperationIdentifier{
		{
			Index: 0,
		},
	}
	baseFeeVaultAddress := common.BaseFeeVault.Hex()
	baseFeeVaultAmount := evmClient.Amount(tx.FeeBurned, sdkTypes.Currency)
	L1FeeVaultRelatedOps := []*RosettaTypes.OperationIdentifier{
		{
			Index: 0,
		},
	}
	L1FeeVaultAddress := common.L1FeeVault.Hex()
	L1FeeVaultAmount := evmClient.Amount(receipt.L1Fee, sdkTypes.Currency)

	ops := []*RosettaTypes.Operation{
		GenerateOp(0, nil, opType, opStatus, fromAddress, fromAmount, nil),
		GenerateOp(1, sequencerRelatedOps, opType, opStatus, sequencerAddress, sequencerAmount, nil),
		GenerateOp(2, baseFeeVaultRelatedOps, opType, opStatus, baseFeeVaultAddress, baseFeeVaultAmount, nil),
		GenerateOp(3, L1FeeVaultRelatedOps, opType, opStatus, L1FeeVaultAddress, L1FeeVaultAmount, nil),
	}

	return ops, nil
}
