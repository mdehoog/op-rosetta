package client

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

// GetNativeTransferGasLimit is Ethereum's custom implementation of estimating gas.
func (c *OpClient) GetNativeTransferGasLimit(ctx context.Context, toAddress string, fromAddress string, value *big.Int) (uint64, error) {
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
