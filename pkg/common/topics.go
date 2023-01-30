package common

import (
	"github.com/ethereum/go-ethereum/core/types"
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

// HashEvent encodes an event string into a topic hash.
func HashEvent(topic string) string {
	return Keccak256(topic)
}
