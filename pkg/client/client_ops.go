package client

import (
	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
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

	feeOps, err := handlers.FeeOps(tx)
	if err != nil {
		return nil, err
	}
	ops = append(ops, feeOps...)

	// TODO(Jingfu): handle burn tx
	if tx.Transaction.IsDepositTx() && !tx.Transaction.IsSystemTx() {
		ops = append(ops, handlers.MintOps(tx, len(ops))...)
	} else {
		ops = append(ops, handlers.TraceOps(tx.Trace, len(ops))...)
	}

	return ops, nil
}
