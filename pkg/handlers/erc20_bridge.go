package handlers

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mdehoog/op-rosetta/pkg/common"
)

// Erc20BridgeOps constructs a list of [RosettaTypes.Operation]s for ERC20 transactions when bridging.
func Erc20BridgeOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	ops := []*RosettaTypes.Operation{}
	ops = append(ops, NativeInitializationOps(tx, startIndex)...)

	return ops
}

// NativeInitializationOps constructs a list of [RosettaTypes.Operation]s for native token initializations.
//
// Native token initializations occur when a native token is being *sent* to a non-native chain.
// In this case, `ERC20BridgeInitiated` and `Transfer` events are emitted FROM the sender TO the bridge contract.
func NativeInitializationOps(tx *evmClient.LoadedTransaction, startIndex int) []*RosettaTypes.Operation {
	// The bridge transaction must be sent to the [common.L2StandardBridge]
	if *tx.Transaction.To() != common.L2StandardBridge {
		return nil
	}

	// Parse tx receipt logs
	var receiptLogs []*EthTypes.Log
	if tx.Receipt != nil {
		receiptLogs = tx.Receipt.Logs
	}

	// Grab the transfer log
	var transferLog *EthTypes.Log
	for _, log := range receiptLogs {
		// Check if this is a transfer event
		if common.ContainsTopic(log, common.TransferEvent) {
			transferLog = log
		}
	}
	if transferLog == (&EthTypes.Log{}) || transferLog == nil {
		return nil
	}

	// Return the associated operation
	opIndex := int64(startIndex)
	opType := sdkTypes.OpErc20Transfer
	opStatus := sdkTypes.SuccessStatus
	fromAddress := evmClient.MustChecksum(tx.From.String())
	amount := evmClient.Erc20Amount(transferLog.Data, transferLog.Address, sdkTypes.Currency.Symbol, sdkTypes.Currency.Decimals, false)

	return []*RosettaTypes.Operation{
		GenerateOp(opIndex, nil, opType, opStatus, fromAddress, amount, nil),
	}
}
