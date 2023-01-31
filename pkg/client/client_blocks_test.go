package client_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	client "github.com/mdehoog/op-rosetta/pkg/client"
	mocks "github.com/mdehoog/op-rosetta/pkg/client/mocks"

	SdkClient "github.com/coinbase/rosetta-geth-sdk/client"
	EthCommon "github.com/ethereum/go-ethereum/common"
	EthTypes "github.com/ethereum/go-ethereum/core/types"
	EthRpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ClientBlocksTestSuite is a test suite for the client block handler.
type ClientBlocksTestSuite struct {
	suite.Suite

	internalClient *mocks.InternalClient
	client         *client.OpClient
}

// SetupTest sets up the test suite.
func (testSuite *ClientBlocksTestSuite) SetupTest() {
	testSuite.internalClient = new(mocks.InternalClient)
	testSuite.client = &client.OpClient{testSuite.internalClient}
}

// TestClientBlocks runs the ClientBlocksTestSuite.
func TestClientBlocks(t *testing.T) {
	suite.Run(t, new(ClientBlocksTestSuite))
}

// TestGetBlockReceiptsBatchCallErrors tests fetching block receipts from the op client with failing batch calls.
func (testSuite *ClientBlocksTestSuite) TestGetBlockReceiptsBatchCallErrors() {
	// Construct arguments
	ctx := context.Background()
	hash := EthCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	gasPrice := big.NewInt(10000)
	blockNumber := big.NewInt(1)
	blockNumberString := blockNumber.String()
	to := EthCommon.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	myTx := EthTypes.NewTransaction(
		0,
		to,
		big.NewInt(0),
		0,
		gasPrice,
		nil,
	)
	txs := []SdkClient.RPCTransaction{
		{
			Tx: myTx,
			TxExtraInfo: SdkClient.TxExtraInfo{
				BlockNumber: &blockNumberString,
				BlockHash:   &hash,
				From:        &to,
				TxHash:      &hash,
			},
		},
	}
	baseFee := big.NewInt(10000)

	// Mock the internall call to the mock client
	ethReceipt := EthTypes.Receipt{
		// Consensus fields: These fields are defined by the Yellow Paper
		Type:              0,
		PostState:         []byte{0x00},
		Status:            1,
		CumulativeGasUsed: 0,
		Bloom:             EthTypes.BytesToBloom([]byte{0x00}),
		Logs:              []*EthTypes.Log{},
		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		TxHash:          hash,
		ContractAddress: to,
		GasUsed:         0,
		// transaction corresponding to this receipt.
		BlockHash:        hash,
		BlockNumber:      blockNumber,
		TransactionIndex: 0,
		// OVM legacy: extend receipts with their L1 price (if a rollup tx)
		// IGNORED
	}
	testSuite.internalClient.On("BatchCallContext", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).([]EthRpc.BatchElem)
		arg[0].Result = ethReceipt
		arg[0].Error = fmt.Errorf("error")
	})

	// Execute and validate the call
	_, err := testSuite.client.GetBlockReceipts(ctx, hash, txs, baseFee)
	testSuite.Equal(fmt.Errorf("error"), err)
	// testSuite.Equal([]*SdkClient.RosettaTxReceipt{}, txReceipts)
}

// TestGetBlockReceiptsEmptyTxs tests fetching block receipts from the op client with no transactions.
func (testSuite *ClientBlocksTestSuite) TestGetBlockReceiptsEmptyTxs() {
	// Construct arguments
	ctx := context.Background()
	hash := EthCommon.HexToHash("0xb358c6958b1cab722752939cbb92e3fec6b6023de360305910ce80c56c3dad9d")
	txs := []SdkClient.RPCTransaction{}
	baseFee := big.NewInt(10000)

	// Execute and validate the call
	txReceipts, err := testSuite.client.GetBlockReceipts(ctx, hash, txs, baseFee)
	testSuite.NoError(err)
	testSuite.Equal([]*SdkClient.RosettaTxReceipt{}, txReceipts)
}
