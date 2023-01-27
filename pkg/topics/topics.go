package topics

import (
	"github.com/ethereum-optimism/optimism/l2geth/core/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// ContainsTopic tests if a log contains a given topic.
// It is expected that the topic is encoded as a hex string.
// To encode a topic string, use the below [EncodeEventString] function.
func ContainsTopic(log *types.Log, topic string) bool {
	for _, t := range log.Topics {
		hex := t.Hex()
		if hex == topic {
			return true
		}
	}
	return false
}

// EncodeEventString encodes a string into a topic.
func EncodeEventString(topicString string) string {
	keccak := crypto.Keccak256([]byte(topicString))
	encodedTransferMethod := hexutil.Encode(keccak)
	return encodedTransferMethod
}
