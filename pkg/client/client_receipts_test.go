package client_test

import (
	"context"
	"math/big"
	"testing"

	client "github.com/mdehoog/op-rosetta/pkg/client"
	mocks "github.com/mdehoog/op-rosetta/pkg/client/mocks"

	SdkClient "github.com/coinbase/rosetta-geth-sdk/client"
	EthCommon "github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ClientReceiptsTestSuite tests [GetTransactionReceipt].
type ClientReceiptsTestSuite struct {
	suite.Suite

	internalClient *mocks.InternalClient
	client         *client.OpClient
}

// SetupTest sets up the test suite.
func (testSuite *ClientReceiptsTestSuite) SetupTest() {
	testSuite.internalClient = new(mocks.InternalClient)
	testSuite.client = &client.OpClient{testSuite.internalClient}
}

// TestClientReceipts runs the ClientReceiptsTestSuite.
func TestClientReceipts(t *testing.T) {
	suite.Run(t, new(ClientReceiptsTestSuite))
}

// TestGetTransactionReceipt tests fetching a transaction receipt from the op client.
func (testSuite *ClientBlocksTestSuite) TestGetTransactionReceipt() {
	// Construct arguments
	ctx := context.Background()
	hash := EthCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	gasPrice := big.NewInt(10000)
	blockNumber := big.NewInt(1)
	to := EthCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	myTx := EthTypes.NewTransaction(
		0,
		to,
		big.NewInt(0),
		0,
		gasPrice,
		nil,
	)
	loadedTx := SdkClient.LoadedTransaction{
		TxHash:      &hash,
		Transaction: myTx,
	}

	// Mock the internall call to the mock client
	ethReceipt := EthTypes.Receipt{
		Type:              0,
		PostState:         []byte{0x00},
		Status:            1,
		CumulativeGasUsed: 0,
		Bloom:             EthTypes.BytesToBloom([]byte{0x00}),
		Logs:              []*EthTypes.Log{},
		TxHash:            hash,
		ContractAddress:   to,
		GasUsed:           0,
		BlockHash:         hash,
		BlockNumber:       blockNumber,
		TransactionIndex:  0,
	}
	testSuite.internalClient.On("CallContext", ctx, mock.Anything, "eth_getTransactionReceipt", &hash).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(**EthTypes.Receipt)
		*arg = &ethReceipt
	})

	// Build rosetta tx receipt response
	gasPrice, _ = SdkClient.EffectiveGasPrice(loadedTx.Transaction, loadedTx.BaseFee)
	gasUsed := new(big.Int).SetUint64(ethReceipt.GasUsed)
	feeAmount := new(big.Int).Mul(gasUsed, gasPrice)
	if ethReceipt.L1Fee != nil {
		feeAmount.Add(feeAmount, ethReceipt.L1Fee)
	}
	expectedRosettaTxReceipt := SdkClient.RosettaTxReceipt{
		GasPrice:       gasPrice,
		GasUsed:        gasUsed,
		Logs:           ethReceipt.Logs,
		RawMessage:     nil,
		TransactionFee: feeAmount,
	}

	// Execute and validate the call
	rosettaTxReceipt, err := testSuite.client.GetTransactionReceipt(ctx, &loadedTx)
	testSuite.NoError(err)
	testSuite.Equal(expectedRosettaTxReceipt, *rosettaTxReceipt)
}
