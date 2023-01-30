package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
)

func (c *OpClient) GetBlockReceipts(
	ctx context.Context,
	blockHash common.Hash,
	txs []evmClient.RPCTransaction,
	baseFee *big.Int,
) ([]*evmClient.RosettaTxReceipt, error) {
	receipts := make([]*evmClient.RosettaTxReceipt, len(txs))
	if len(txs) == 0 {
		return receipts, nil
	}

	ethReceipts := make([]*EthTypes.Receipt, len(txs))
	reqs := make([]rpc.BatchElem, len(txs))
	for i := range reqs {
		reqs[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{txs[i].TxExtraInfo.TxHash.String()},
			Result: &ethReceipts[i],
		}
	}

	maxBatchSize := 25
	for i := 0; i < len(txs); i += maxBatchSize {
		if i+maxBatchSize < len(txs) {
			if err := c.BatchCallContext(ctx, reqs[i:i+maxBatchSize]); err != nil {
				return nil, err
			}
		} else {
			if err := c.BatchCallContext(ctx, reqs[i:]); err != nil {
				return nil, err
			}
		}
	}

	for i := range reqs {
		if reqs[i].Error != nil {
			return nil, reqs[i].Error
		}

		gasPrice, err := evmClient.EffectiveGasPrice(txs[i].Tx, baseFee)
		if err != nil {
			return nil, err
		}
		gasUsed := new(big.Int).SetUint64(ethReceipts[i].GasUsed)
		feeAmount := new(big.Int).Mul(gasUsed, gasPrice)
		if ethReceipts[i].L1Fee != nil {
			feeAmount.Add(feeAmount, ethReceipts[i].L1Fee)
		}

		receiptJSON, err := ethReceipts[i].MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("unable to marshal receipt for %x: %v", txs[i].Tx.Hash().Hex(), err)
		}
		receipt := &evmClient.RosettaTxReceipt{
			Type:     ethReceipts[i].Type,
			GasPrice: gasPrice,
			GasUsed:  gasUsed,
			Logs:     ethReceipts[i].Logs,
			// This is a hack to get around the fact that the RosettaTxReceipt doesn't contain L1 fees. We add the raw receipt here so we can access other rollup fields later
			RawMessage:     receiptJSON,
			TransactionFee: feeAmount,
		}

		receipts[i] = receipt

		if ethReceipts[i] == nil {
			return nil, fmt.Errorf("got empty receipt for %x", txs[i].Tx.Hash().Hex())
		}

		if ethReceipts[i].BlockHash != blockHash {
			return nil, fmt.Errorf(
				"%w: expected block hash %s for Transaction but got %s",
				sdkTypes.ErrClientBlockOrphaned,
				blockHash.Hex(),
				ethReceipts[i].BlockHash.Hex(),
			)
		}
	}

	return receipts, nil
}
