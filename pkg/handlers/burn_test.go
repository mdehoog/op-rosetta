package handlers_test

import (
	"math/big"
	"testing"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	common "github.com/mdehoog/op-rosetta/pkg/common"
	handlers "github.com/mdehoog/op-rosetta/pkg/handlers"
	suite "github.com/stretchr/testify/suite"
)

type BurnTestSuite struct {
	suite.Suite
}

// TestBurn runs the BurnTestSuite.
func TestBurn(t *testing.T) {
	suite.Run(t, new(BurnTestSuite))
}

// TestInvalidDestination tests that a [evmClient.LoadedTransaction] with an invalid destination is not a burn.
func (testSuite *BurnTestSuite) TestInvalidDestination() {
	// Construct a loaded transaction with an invalid burn destination (not L2ToL1MessagePasser).
	h := gethCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	gasPrice := big.NewInt(10000)
	myTx := EthTypes.NewTransaction(
		0,
		gethCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0),
		0,
		gasPrice,
		nil,
	)
	loadedTxn := &evmClient.LoadedTransaction{
		TxHash:      &h,
		Transaction: myTx,
	}

	// BurnOps should return nil for a transaction with an invalid destination.
	ops := handlers.BurnOps(loadedTxn, 0)
	testSuite.Nil(ops)
}

// TestValidBurn tests that [handlers.BurnOps] correctly constructs [RosettaTypes.Operation],
// given a [evmClient.LoadedTransaction] with the correct destination address.
func (testSuite *BurnTestSuite) TestValidBurn() {
	// Construct a loaded transaction with the correct burn destination (L2ToL1MessagePasser).
	// Note: this hash is incorrect and was hijacked from the above transaction.
	h := gethCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	gasPrice := big.NewInt(10000)
	amount := big.NewInt(100)
	index := 1
	myTx := EthTypes.NewTransaction(
		0,
		common.L2ToL1MessagePasser,
		amount,
		0,
		gasPrice,
		nil,
	)
	from := gethCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	loadedTxn := &evmClient.LoadedTransaction{
		From:        &from,
		TxHash:      &h,
		Transaction: myTx,
	}

	// BurnOps should successfully construct a BurnOp for a transaction with the correct destination.
	ops := handlers.BurnOps(loadedTxn, index)
	testSuite.Equal(ops, []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(index),
			},
			Type:   common.BurnOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: from.String(),
			},
			Amount: evmClient.Amount(amount, sdkTypes.Currency),
		},
	})
}
