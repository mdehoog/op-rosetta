package config

import (
	"github.com/coinbase/rosetta-geth-sdk/configuration"
)

const (
	// DefaultBlockchain is Optimism.
	DefaultBlockchain string = "Optimism"

	// BlockchainEnv is the environment variable
	// read to determine the blockchain.
	BlockchainEnv = "BLOCKCHAIN"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "ETH"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 18

	// Online is when the implementation is permitted
	// to make outbound connections.
	Online configuration.Mode = "ONLINE"

	// Offline is when the implementation is not permitted
	// to make outbound connections.
	Offline configuration.Mode = "OFFLINE"

	// ModeEnv is the environment variable read
	// to determine mode.
	ModeEnv = "MODE"

	// NetworkEnv is the environment variable
	// read to determine network.
	NetworkEnv = "NETWORK"

	// TransitionBlockHashEnv is the environment variable
	// from which to read the transition block hash.
	TransitionBlockHashEnv = "TRANSITION_BLOCK_HASH"

	// ChainConfigEnv is the environment variable from
	// which to read the chain configuration, defined as
	// JSON (or pointing to a JSON file).
	ChainConfigEnv = "CHAIN_CONFIG"

	// PortEnv is the environment variable
	// read to determine the port for the Rosetta
	// implementation.
	PortEnv = "PORT"

	// FilterTokensEnv is the environment variable
	// read to determine if we will filter tokens
	// using our token white list
	FilterTokensEnv = "FILTER"

	// TokenListEnv is the environment variable
	// from which to read the list of tokens, defined
	// as JSON (or pointing to a JSON file).
	TokenListEnv = "TOKENS"

	// GethEnv is an optional environment variable
	// used to connect rosetta-ethereum to an already
	// running geth node.
	GethEnv = "GETH"

	// DefaultGethURL is the default URL for
	// a running geth node. This is used
	// when GethEnv is not populated.
	DefaultGethURL = "http://127.0.0.1:8545"

	// SkipGethAdminEnv is an optional environment variable
	// to skip geth `admin` calls which are typically not supported
	// by hosted node services. When not set, defaults to false.
	SkipGethAdminEnv = "SKIP_GETH_ADMIN"

	// GenesisBlockIndex is the index of the genesis block.
	GenesisBlockIndex = int64(0)
)
