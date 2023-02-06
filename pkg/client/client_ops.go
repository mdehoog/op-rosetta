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
		trace := tx.Trace[0]
		traceType := strings.ToUpper(trace.Type)
		opStatus := sdkTypes.SuccessStatus
		from := evmClient.MustChecksum(trace.From.String())
		to := evmClient.MustChecksum(trace.To.String())
		metadata := map[string]interface{}{}

		if from != to {
			feeOps, err := handlers.FeeOps(tx)
			if err != nil {
				return nil, err
			}
			ops = append(ops, feeOps...)
		}

		toOp := &RosettaTypes.Operation{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(len(ops) + 0),
			},
			Type:   traceType,
			Status: RosettaTypes.String(opStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: to,
			},
			Amount: &RosettaTypes.Amount{
				Value:    trace.Value.String(),
				Currency: sdkTypes.Currency,
			},
			Metadata: metadata,
		}
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
