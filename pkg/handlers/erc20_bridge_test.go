package handlers_test

import (
	"testing"

	suite "github.com/stretchr/testify/suite"
)

type Erc20BridgeTestSuite struct {
	suite.Suite
}

// TestErc20Bridging runs the Erc20BridgeTestSuite.
func TestErc20Bridging(t *testing.T) {
	suite.Run(t, new(Erc20BridgeTestSuite))
}

// TestValidNativeInitialization tests that a valid ERC20 transfer was made to the L2StandardBridge contract.
// This is when a native token is being bridged from L2 to L1.
func (testSuite *Erc20BridgeTestSuite) TestValidNativeInitialization() {
	// TODO:
}
