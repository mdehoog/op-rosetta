package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/params"

	"github.com/coinbase/rosetta-geth-sdk/configuration"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
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

	// GenesisBlockHashEnv is the environment variable
	// from which to read the genesis block hash.
	GenesisBlockHashEnv = "GENESIS_BLOCK_HASH"

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

// LoadConfiguration attempts to create a new Configuration
// using the ENVs in the environment.
func LoadConfiguration() (*configuration.Configuration, error) {
	config := &configuration.Configuration{}

	mode := os.Getenv(ModeEnv)
	modeValue := configuration.Mode(mode)

	switch modeValue {
	case Online:
		config.Mode = Online
	case Offline:
		config.Mode = Offline
	case "":
		return nil, fmt.Errorf("%s must be populated", ModeEnv)
	default:
		return nil, fmt.Errorf("%s is not a valid mode", modeValue)
	}

	blockchain := os.Getenv(BlockchainEnv)
	if blockchain == "" {
		blockchain = DefaultBlockchain
	}

	networkValue := os.Getenv(NetworkEnv)
	genesisBlockHash := &RosettaTypes.BlockIdentifier{
		Index: GenesisBlockIndex,
		Hash:  os.Getenv(GenesisBlockHashEnv),
	}

	chainConfigJson := os.Getenv(ChainConfigEnv)
	if file, err := os.ReadFile(chainConfigJson); err == nil {
		// if the envvar points to a file, read it; otherwise the envvar contents is expected to be JSON
		chainConfigJson = string(file)
	}
	if chainConfigJson == "" {
		return nil, fmt.Errorf("%s not set", ChainConfigEnv)
	}

	chainConfig := &params.ChainConfig{}
	err := json.Unmarshal([]byte(chainConfigJson), &chainConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to parse chain config: %w", err)
	}

	config.Network = &RosettaTypes.NetworkIdentifier{
		Blockchain: blockchain,
		Network:    networkValue,
	}
	config.GenesisBlockIdentifier = genesisBlockHash
	config.ChainConfig = chainConfig

	config.GethURL = DefaultGethURL
	envGethURL := os.Getenv(GethEnv)
	if len(envGethURL) > 0 {
		config.RemoteGeth = true
		config.GethURL = envGethURL
	}

	config.SkipGethAdmin = false
	envSkipGethAdmin := os.Getenv(SkipGethAdminEnv)
	if len(envSkipGethAdmin) > 0 {
		val, err := strconv.ParseBool(envSkipGethAdmin)
		if err != nil {
			return nil, fmt.Errorf("unable to parse SKIP_GETH_ADMIN %s: %w", envSkipGethAdmin, err)
		}
		config.SkipGethAdmin = val
	}

	portValue := os.Getenv(PortEnv)
	if len(portValue) == 0 {
		return nil, errors.New("PORT must be populated")
	}

	port, err := strconv.Atoi(portValue)
	if err != nil || len(portValue) == 0 || port <= 0 {
		return nil, fmt.Errorf("unable to parse port %s: %w", portValue, err)
	}
	config.Port = port

	filterTokens := true
	filterTokensStr := os.Getenv(FilterTokensEnv)
	if len(filterTokensStr) > 0 {
		filterTokens, err = strconv.ParseBool(filterTokensStr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse token filter %s: %w", filterTokensStr, err)
		}
	}

	tokenListJson := os.Getenv(TokenListEnv)
	if file, err := os.ReadFile(tokenListJson); err == nil {
		// if the envvar points to a file, read it; otherwise the envvar contents is expected to be JSON
		tokenListJson = string(file)
	}
	if tokenListJson == "" {
		tokenListJson = "[]"
	}

	var payload []configuration.Token
	err = json.Unmarshal([]byte(tokenListJson), &payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", tokenListJson, err)
	}

	config.RosettaCfg = configuration.RosettaConfig{
		SupportRewardTx: false,
		TraceType:       configuration.GethNativeTrace,
		Currency: &RosettaTypes.Currency{
			Symbol:   Symbol,
			Decimals: Decimals,
		},
		TracePrefix:     "optrace",
		FilterTokens:    filterTokens,
		TokenWhiteList:  payload,
		SupportsSyncing: true,
	}

	return config, nil
}
