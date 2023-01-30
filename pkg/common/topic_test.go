package common_test

import (
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mdehoog/op-rosetta/pkg/common"

	"github.com/stretchr/testify/suite"
)

type TopicTestSuite struct {
	suite.Suite
}

// TestTopics runs the TopicTestSuite.
func TestTopics(t *testing.T) {
	suite.Run(t, new(TopicTestSuite))
}

// ERC20_TRANSFER_EVENT_HASH is the keccak256 hash of the Transfer(address,address,uint256) event.
const ERC20_TRANSFER_EVENT_HASH = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// TestHashEvent tests encoding an event topic string.
func (testSuite *TopicTestSuite) TestHashEvent() {
	topicString := "Transfer(address,address,uint256)"
	encoded := common.Keccak256(topicString)
	testSuite.Equal(ERC20_TRANSFER_EVENT_HASH, encoded)
}

// TestContainsTopic tests if a log contains a given topic.
func (testSuite *TopicTestSuite) TestContainsTopic() {
	topicString := "Transfer(address,address,uint256)"
	keccak := crypto.Keccak256([]byte(topicString))
	encodedTransferMethod := hexutil.Encode(keccak)
	testSuite.Equal(ERC20_TRANSFER_EVENT_HASH, encodedTransferMethod)
	l := &types.Log{Topics: []gethCommon.Hash{gethCommon.HexToHash(ERC20_TRANSFER_EVENT_HASH)}}
	common.ContainsTopic(l, encodedTransferMethod)
}
