package client_test

import (
	"context"
	"math/big"
	"testing"

	client "github.com/mdehoog/op-rosetta/pkg/client"
	mocks "github.com/mdehoog/op-rosetta/pkg/client/mocks"

	Ethereum "github.com/ethereum/go-ethereum"
	EthCommon "github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/suite"
)

// ClientGasTestSuite is a test suite for the client block handler.
type ClientGasTestSuite struct {
	suite.Suite

	internalClient *mocks.InternalClient
	client         *client.OpClient
}

// SetupTest sets up the test suite.
func (testSuite *ClientGasTestSuite) SetupTest() {
	testSuite.internalClient = new(mocks.InternalClient)
	testSuite.client = &client.OpClient{testSuite.internalClient}
}

// TestClientGas runs the ClientGasTestSuite.
func TestClientGas(t *testing.T) {
	suite.Run(t, new(ClientGasTestSuite))
}

// TestGetNativeTransferGasLimit tests fetching the transfer gas limit from the op client.
func (testSuite *ClientGasTestSuite) TestGetNativeTransferGasLimit() {
	// construct args
	// vb1 -> vb2 :P
	toAddress := EthCommon.HexToAddress("5763A0B54AA391d7B2255b2F824cef80E8F343d0")
	fromAddress := EthCommon.HexToAddress("234D953a9404Bf9DbC3b526271d440cD2870bCd2")
	value := big.NewInt(10000)
	ctx := context.Background()

	// Mock the internal EstimateGas call
	mockGasLimit := uint64(10000)
	testSuite.internalClient.On("EstimateGas", ctx, Ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: value,
	}).Return(mockGasLimit, nil)

	// Execute and validate the call
	gas, err := testSuite.client.GetNativeTransferGasLimit(ctx, toAddress.Hex(), fromAddress.Hex(), value)
	testSuite.NoError(err)
	testSuite.Equal(mockGasLimit, gas)
}

// TestGasInvalidToAddress tests calling [GetNativeTransferGasLimit] with an invalid `toAddress` argument.
func (testSuite *ClientGasTestSuite) TestGasInvalidToAddress() {
	toAddress := ""
	fromAddress := EthCommon.HexToAddress("234D953a9404Bf9DbC3b526271d440cD2870bCd2")
	value := big.NewInt(10000)
	ctx := context.Background()
	_, err := testSuite.client.GetNativeTransferGasLimit(ctx, toAddress, fromAddress.Hex(), value)
	testSuite.Nil(err)
}

// TestGasInvalidValue tests calling [GetNativeTransferGasLimit] with an invalid `value` argument.
func (testSuite *ClientGasTestSuite) TestGasInvalidValue() {
	toAddress := EthCommon.HexToAddress("5763A0B54AA391d7B2255b2F824cef80E8F343d0")
	fromAddress := EthCommon.HexToAddress("234D953a9404Bf9DbC3b526271d440cD2870bCd2")
	ctx := context.Background()
	_, err := testSuite.client.GetNativeTransferGasLimit(ctx, toAddress.Hex(), fromAddress.Hex(), nil)
	testSuite.Nil(err)
}
