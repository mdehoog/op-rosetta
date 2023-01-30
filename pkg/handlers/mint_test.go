package handlers_test

import (
	"math/big"
	"testing"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	sdkTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	handlers "github.com/mdehoog/op-rosetta/pkg/handlers"
	suite "github.com/stretchr/testify/suite"
)

type MintTestSuite struct {
	suite.Suite
}

// TestMint runs the MintTestSuite.
func TestMint(t *testing.T) {
	suite.Run(t, new(MintTestSuite))
}

// TestInvalidDeposit tests that a non-deposit [evmClient.LoadedTransaction] is not handled by MintOps.
func (testSuite *MintTestSuite) TestInvalidDeposit() {
	// Construct a random transaction (non-DepositTx)
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

	// MintOps should return nil for non-min transactions.
	ops := handlers.MintOps(loadedTxn, 0)
	testSuite.Nil(ops)
}

// TestValidMint tests that [handlers.MintOps] correctly constructs [RosettaTypes.Operation],
// given a [evmClient.LoadedTransaction] with a Mint transaction.
func (testSuite *MintTestSuite) TestValidMint() {
	// Construct a loaded mint transaction.
	// Note: this hash is incorrect and was hijacked from the above transaction.
	h := gethCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	gasPrice := big.NewInt(10000)
	amount := big.NewInt(100)
	index := 1
	to := gethCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	myTx := EthTypes.DepositTx{
		From:                gethCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		Value:               amount,
		To:                  &to,
		Mint:                amount,
		Gas:                 gasPrice.Uint64(),
		IsSystemTransaction: false,
	}
	from := gethCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	loadedTxn := &evmClient.LoadedTransaction{
		From:        &from,
		TxHash:      &h,
		Transaction: EthTypes.NewTx(&myTx),
	}

	// MintOps should successfully construct a Mint operation.
	ops := handlers.MintOps(loadedTxn, index)
	testSuite.Equal(ops, []*RosettaTypes.Operation{
		{
			OperationIdentifier: &RosettaTypes.OperationIdentifier{
				Index: int64(index),
			},
			Type:   handlers.MintOpType,
			Status: RosettaTypes.String(sdkTypes.SuccessStatus),
			Account: &RosettaTypes.AccountIdentifier{
				Address: from.String(),
			},
			Amount: evmClient.Amount(amount, sdkTypes.Currency),
		},
	})
}
