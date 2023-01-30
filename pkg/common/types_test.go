package common_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mdehoog/op-rosetta/pkg/common"
	suite "github.com/stretchr/testify/suite"
)

type CommonTestSuite struct {
	suite.Suite
}

// TestCommon runs the CommonTestSuite.
func TestCommon(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func hash(s string) string {
	keccak := crypto.Keccak256([]byte(s))
	return hexutil.Encode(keccak)
}

// TestEventTopics tests that the event strings are correctly hashed to their topics.
func (testSuite *CommonTestSuite) TestEventTopics() {
	burnEventHash := "0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5"
	mintEventHash := "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"
	transferEventHash := "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	erc20BridgeInitiatedEventHash := "0x7ff126db8024424bbfd9826e8ab82ff59136289ea440b04b39a0df1b03b9cabf"
	erc20BridgeFinalizedEventHash := "0xd59c65b35445225835c83f50b6ede06a7be047d22e357073e250d9af537518cd"
	testSuite.Equal(hash(common.BurnEvent), burnEventHash)
	testSuite.Equal(hash(common.MintEvent), mintEventHash)
	testSuite.Equal(hash(common.TransferEvent), transferEventHash)
	testSuite.Equal(hash(common.ERC20BridgeInitiatedEvent), erc20BridgeInitiatedEventHash)
	testSuite.Equal(hash(common.ERC20BridgeFinalizedEvent), erc20BridgeFinalizedEventHash)
}
