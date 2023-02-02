package config

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

// LoadConfiguration attempts to create a new [configuration.Configuration]
// using environment variables.
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

	// Graceful tokenlist unmarshalling
	tokenListJsonFilename := os.Getenv(TokenListEnv)
	tokenListJsonFile, err := ReadTokenConfig(tokenListJsonFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token list %s: %w", tokenListJsonFilename, err)
	}
	tokenWhiteList, err := UnmarshalTokenConfig([]byte(tokenListJsonFile), *chainConfig.ChainID)
	if err != nil {
		return nil, err
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
		TokenWhiteList:  tokenWhiteList,
		SupportsSyncing: true,
	}

	return config, nil
}
