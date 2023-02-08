package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
)

// GetTransactionReceipt returns the receipt for an [evmClient.LoadedTransaction].
func (c *OpClient) GetTransactionReceipt(
	ctx context.Context,
	tx *evmClient.LoadedTransaction,
) (*evmClient.RosettaTxReceipt, error) {
	var r *EthTypes.Receipt
	err := c.CallContext(ctx, &r, "eth_getTransactionReceipt", tx.TxHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	gasPrice, err := evmClient.EffectiveGasPrice(tx.Transaction, tx.BaseFee)
	if err != nil {
		return nil, err
	}
	gasUsed := new(big.Int).SetUint64(r.GasUsed)
	feeAmount := new(big.Int).Mul(gasUsed, gasPrice)
	if r.L1Fee != nil {
		feeAmount.Add(feeAmount, r.L1Fee)
	}

	return &evmClient.RosettaTxReceipt{
		GasPrice:       gasPrice,
		GasUsed:        gasUsed,
		Logs:           r.Logs,
		RawMessage:     nil,
		TransactionFee: feeAmount,
	}, err
}
