package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/coinbase/rosetta-geth-sdk/configuration"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	"github.com/coinbase/rosetta-geth-sdk/services"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
)

type OpClient struct {
	*evmClient.SDKClient
}

const (
	L1ToL2DepositType = 126
)

func (c *OpClient) ParseOps(
	tx *evmClient.LoadedTransaction,
) ([]*types.Operation, error) {
	var ops []*types.Operation

	if tx.Receipt.Type == L1ToL2DepositType && len(tx.Trace) > 0 {
		trace := tx.Trace[0]
		traceType := strings.ToUpper(trace.Type)
		opStatus := sdkTypes.SuccessStatus
		from := evmClient.MustChecksum(trace.From.String())
		to := evmClient.MustChecksum(trace.To.String())
		metadata := map[string]interface{}{}

		if from != to {
			feeOps := services.FeeOps(tx)
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

	feeOps := services.FeeOps(tx)
	ops = append(ops, feeOps...)

	traceOps := services.TraceOps(tx.Trace, len(ops))
	ops = append(ops, traceOps...)

	return ops, nil
}

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
	if err := c.BatchCallContext(ctx, reqs); err != nil {
		return nil, err
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

		receipt := &evmClient.RosettaTxReceipt{
			Type: ethReceipts[i].Type,
			GasPrice:       gasPrice,
			GasUsed:        gasUsed,
			Logs:           ethReceipts[i].Logs,
			RawMessage:     nil,
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

	return &evmClient.RosettaTxReceipt{
		GasPrice:       gasPrice,
		GasUsed:        gasUsed,
		Logs:           r.Logs,
		RawMessage:     nil,
		TransactionFee: feeAmount,
	}, err
}

// GetNativeTransferGasLimit is Ethereum's custom implementation of estimating gas.
func (c *OpClient) GetNativeTransferGasLimit(ctx context.Context, toAddress string,
	fromAddress string, value *big.Int) (uint64, error) {
	if len(toAddress) == 0 || value == nil {
		// We guard against malformed inputs that may have been generated using
		// a previous version of asset's rosetta
		return 21000, nil
	}
	to := common.HexToAddress(toAddress)
	return c.EstimateGas(ctx, ethereum.CallMsg{
		From:  common.HexToAddress(fromAddress),
		To:    &to,
		Value: value,
	})
}

// NewOpClient creates a client that can interact with the Optimism network.
func NewOpClient(cfg *configuration.Configuration) (*OpClient, error) {
	client, err := evmClient.NewClient(cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize client: %w", err)
	}
	return &OpClient{client}, nil
}
