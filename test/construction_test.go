package test

import (
	"context"
	"encoding/json"
	"testing"

	config "github.com/mdehoog/op-rosetta/pkg/config"
	utils "github.com/mdehoog/op-rosetta/pkg/utils"

	RosettaConfiguration "github.com/coinbase/rosetta-geth-sdk/configuration"
	RosettaConstruction "github.com/coinbase/rosetta-geth-sdk/services/construction"
	RosettaGethTypes "github.com/coinbase/rosetta-geth-sdk/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	EthParams "github.com/ethereum/go-ethereum/params"
	assert "github.com/stretchr/testify/assert"
)

func TestPreprocessERC20(t *testing.T) {
	networkID := &RosettaTypes.NetworkIdentifier{
		Blockchain: config.DefaultBlockchain,
		Network:    "testnet",
	}
	cfg := &RosettaConfiguration.Configuration{
		Mode:        RosettaConfiguration.ModeOnline,
		Network:     networkID,
		ChainConfig: EthParams.TestChainConfig,
	}
	servicer := RosettaConstruction.NewAPIService(cfg, utils.LoadTypes(), RosettaGethTypes.Errors, nil)
	ctx := context.Background()

	intent := `
	[
		{
			"operation_identifier": {
				"index":0
			},
			"type":"PAYMENT",
			"account": {
				"address":"0x9670d6977d0b10130E5d4916c9134363281B6B0e"
			},
			"amount": {
				"value":"-100000000000",
				"currency":{
					"symbol":"OP",
					"decimals":18,
					"metadata":{
						"token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"
					}
				}
			}
		},
		{
			"operation_identifier": {
				"index":1
			},
			"type":"PAYMENT",
			"account": {
				"address":"0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674"
			},
			"amount": {
				"value":"100000000000",
				"currency": {
					"symbol":"OP",
					"decimals":18,
					"metadata": {
						"token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"
					}
				}
			}
		}
	]
	`
	var ops []*RosettaTypes.Operation
	assert.NoError(t, json.Unmarshal([]byte(intent), &ops))
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&RosettaTypes.ConstructionPreprocessRequest{
			NetworkIdentifier: networkID,
			Operations:        ops,
		},
	)
	assert.Nil(t, err)
	options := map[string]interface{}{
		"from":  "0x9670d6977d0b10130E5d4916c9134363281B6B0e",
		"to":    "0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674",
		"value": "0x0",
		"currency": map[string]interface{}{
			"decimals": float64(18),
			"symbol":   "OP",
			"metadata": map[string]interface{}{
				"contractAddress": "0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A",
			},
		},
	}
	assert.Equal(t, &RosettaTypes.ConstructionPreprocessResponse{
		Options: options,
	}, preprocessResponse)
}

func TestPreprocessGovernanceDelegate(t *testing.T) {
	networkID := &RosettaTypes.NetworkIdentifier{
		Blockchain: config.DefaultBlockchain,
		Network:    "testnet",
	}
	cfg := &RosettaConfiguration.Configuration{
		Mode:        RosettaConfiguration.ModeOnline,
		Network:     networkID,
		ChainConfig: EthParams.TestChainConfig,
	}
	servicer := RosettaConstruction.NewAPIService(cfg, utils.LoadTypes(), RosettaGethTypes.Errors, nil)
	ctx := context.Background()

	intent := `
	[
		{
			"operation_identifier": {
				"index": 0
			},
			"type": "DELEGATE_VOTES",
			"account": {
				"address": "0x9670d6977d0b10130E5d4916c9134363281B6B0e"
			},
			"amount": {
				"value": "0",
				"currency": {
					"symbol": "OP",
					"decimals": 18,
					"metadata": {
						"token_address": "0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"
					}
				}
			}
		},
		{
			"operation_identifier": {
				"index": 1
			},
			"type": "DELEGATE_VOTES",
			"account": {
				"address": "0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674"
			},
			"amount": {
				"value": "0",
				"currency": {
					"symbol": "OP",
					"decimals": 18,
					"metadata": {
					"token_address": "0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"
					}
				}
			}
		}
	]
	`

	var ops []*RosettaTypes.Operation
	assert.NoError(t, json.Unmarshal([]byte(intent), &ops))
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&RosettaTypes.ConstructionPreprocessRequest{
			NetworkIdentifier: networkID,
			Operations:        ops,
		},
	)
	assert.Nil(t, err)
	methodArgs := []interface{}{"0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674"}
	options := map[string]interface{}{
		"from":             "0x9670d6977d0b10130E5d4916c9134363281B6B0e",
		"to":               "0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674",
		"value":            "0x0",
		"contract_address": "0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A",
		"data":             "0x5c19a95c000000000000000000000000705f9aE78b11a3ED5080c053Fa4Fa0c52359c674",
		"method_signature": "delegate(address)",
		"method_args":      methodArgs,
		"currency": map[string]interface{}{
			"decimals": float64(18),
			"symbol":   "OP",
		},
	}
	assert.Equal(t, &RosettaTypes.ConstructionPreprocessResponse{
		Options: options,
	}, preprocessResponse)
}
