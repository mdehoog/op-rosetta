package client

import (
	"strings"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/mdehoog/op-rosetta/pkg/handlers"
)

// ParseOps must be implemented by downstream clients from the [evmClient.SDKClient].
// This method is called by the [evmClient.SDKClient] when populating Block Transactions.
// See [evmClient.SDKClient.PopulateTransactions] for usage.
func (c *OpClient) ParseOps(
	tx *evmClient.LoadedTransaction,
) ([]*RosettaTypes.Operation, error) {
	var ops []*RosettaTypes.Operation

	if tx.Receipt.Type == L1ToL2DepositType && len(tx.Trace) > 0 && tx.Transaction.IsSystemTx() {
		call := tx.Trace[0]
		fromAddress := evmClient.MustChecksum(call.From.String())
		toAddress := evmClient.MustChecksum(call.To.String())

		if fromAddress != toAddress {
			feeOps, err := handlers.FeeOps(tx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, feeOps...)
		}

		opIndex := int64(len(ops) + 0)
		opType := strings.ToUpper(call.Type)
		opStatus := sdkTypes.SuccessStatus
		toAmount := evmClient.Amount(call.Value, sdkTypes.Currency)

		toOp := handlers.GenerateOp(opIndex, nil, opType, opStatus, toAddress, toAmount, nil)
		ops = append(ops, toOp)
		return ops, nil
	}

	feeOps, err := handlers.FeeOps(tx)
	if err != nil {
		return nil, err
	}
	ops = append(ops, feeOps...)

	ops = append(ops, handlers.MintOps(tx, len(ops))...)
	// ops = append(ops, handlers.BurnOps(tx, len(ops))...)
	ops = append(ops, handlers.TraceOps(tx.Trace, len(ops))...)

	return ops, nil
}
