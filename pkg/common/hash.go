package common

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Keccak256 computes the keccak256 hash of the given string.
// Encodes the resulting hash as a hex string with a 0x prefix.
func Keccak256(s string) string {
	keccak := crypto.Keccak256([]byte(s))
	return hexutil.Encode(keccak)
}
