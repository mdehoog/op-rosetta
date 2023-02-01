package config_test

import (
	"testing"

	SdkConfiguration "github.com/coinbase/rosetta-geth-sdk/configuration"
	"github.com/mdehoog/op-rosetta/pkg/config"

	"github.com/stretchr/testify/suite"
)

// TokensTestSuite is a test suite for tokens utilities.
type TokensTestSuite struct {
	suite.Suite
}

// TestTokens runs the TokensTestSuite.
func TestTokens(t *testing.T) {
	suite.Run(t, new(TokensTestSuite))
}

// TestUnmarshalOldTokenConfig tests unmarshalling the outdated token config json format.
func (testSuite *TokensTestSuite) TestUnmarshalOldTokenConfig() {
	contents := `{
		"Mainnet" : {
			"0x4200000000000000000000000000000000000042": true,
			"0x7f5c764cbc14f9669b88837ca1490cca17c31607": true
		},
		"Testnet" : {
			"0xda10009cbd5d07dd0cecc66161fc93d7c9000da1": true
		},
		"Goerli" : {
			"0x7f5c764cbc14f9669b88837ca1490cca17c31607": true
		}
	}`
	c, err := config.UnmarshalTokenConfig(contents)
	testSuite.NoError(err)
	testSuite.Equal([]SdkConfiguration.Token{
		{
			ChainID: 1,
			Address: "0x4200000000000000000000000000000000000042",
		},
		{
			ChainID: 1,
			Address: "0x7f5c764cbc14f9669b88837ca1490cca17c31607",
		},
		{
			ChainID: 420,
			Address: "0xda10009cbd5d07dd0cecc66161fc93d7c9000da1",
		},
		{
			ChainID: 420,
			Address: "0x7f5c764cbc14f9669b88837ca1490cca17c31607",
		},
	}, c)
}

// TestUnmarshalOldTokenConfigExcluded tests unmarshalling the outdated token config json format with excluded tokens.
func (testSuite *TokensTestSuite) TestUnmarshalOldTokenConfigExcluded() {
	contents := `{
		"Mainnet" : {
			"0x4200000000000000000000000000000000000042": true,
			"0x7f5c764cbc14f9669b88837ca1490cca17c31607": false
		},
		"Testnet" : {
			"0xda10009cbd5d07dd0cecc66161fc93d7c9000da1": true
		},
		"Goerli" : {
			"0x7f5c764cbc14f9669b88837ca1490cca17c31607": false
		}
	}`
	c, err := config.UnmarshalTokenConfig(contents)
	testSuite.NoError(err)
	testSuite.Equal([]SdkConfiguration.Token{
		{
			ChainID: 1,
			Address: "0x4200000000000000000000000000000000000042",
		},
		{
			ChainID: 420,
			Address: "0xda10009cbd5d07dd0cecc66161fc93d7c9000da1",
		},
	}, c)
}

// TestUnmarshalTokenConfig tests unmarshalling the token config json format.
func (testSuite *TokensTestSuite) TestUnmarshalTokenConfig() {
	contents := `[
		{
			"chainId": 1,
			"address": "0x7f5c764cbc14f9669b88837ca1490cca17c31607",
			"name": "USD Coin",
			"symbol": "USDC",
			"decimals": 6
		},
		{
			"chainId": 1,
			"address": "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
			"name": "USDT",
			"symbol": "USDT",
			"decimals": 6
		}
	]`
	c, err := config.UnmarshalTokenConfig(contents)
	testSuite.NoError(err)
	testSuite.Equal([]SdkConfiguration.Token{
		{
			ChainID:  1,
			Address:  "0x7f5c764cbc14f9669b88837ca1490cca17c31607",
			Name:     "USD Coin",
			Symbol:   "USDC",
			Decimals: 6,
		},
		{
			ChainID:  1,
			Address:  "0x94b008aa00579c1307b0ef2c499ad98a8ce58e58",
			Name:     "USDT",
			Symbol:   "USDT",
			Decimals: 6,
		},
	}, c)
}
