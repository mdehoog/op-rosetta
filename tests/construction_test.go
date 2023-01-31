package tests

import (
	"context"
	"encoding/json"
	"testing"

	rosettaGethClient "github.com/coinbase/rosetta-geth-sdk/client"
	"github.com/coinbase/rosetta-geth-sdk/configuration"
	construction "github.com/coinbase/rosetta-geth-sdk/services/construction"
	rosettaGethTypes "github.com/coinbase/rosetta-geth-sdk/types"
	rosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/mdehoog/op-rosetta/pkg/config"
	"github.com/mdehoog/op-rosetta/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var (
	OPCurrency = &rosettaTypes.Currency{
		Symbol:   "OP",
		Decimals: 18,
	}
)

func TestPreprocessERC20(t *testing.T) {
	networkID := &rosettaTypes.NetworkIdentifier{
		Blockchain: config.DefaultBlockchain,
		Network:    "testnet",
	}
	cfg := &configuration.Configuration{
		Mode:        configuration.ModeOnline,
		Network:     networkID,
		ChainConfig: params.TestChainConfig,
	}
	servicer := construction.NewAPIService(cfg, utils.LoadTypes(), rosettaGethTypes.Errors, nil)
	ctx := context.Background()

	intent := `
	[
		{"operation_identifier":{"index":0},"type":"PAYMENT","account":{"address":"0x9670d6977d0b10130E5d4916c9134363281B6B0e"},"amount":{"value":"-100000000000","currency":{"symbol":"OP","decimals":18,"metadata":{"token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"}}}},
		{"operation_identifier":{"index":1},"type":"PAYMENT","account":{"address":"0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674"},"amount":{"value":"100000000000","currency":{"symbol":"OP","decimals":18,"metadata":{"token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A"}}}}
	]
	`
	var ops []*rosettaTypes.Operation
	assert.NoError(t, json.Unmarshal([]byte(intent), &ops))
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&rosettaTypes.ConstructionPreprocessRequest{
			NetworkIdentifier: networkID,
			Operations:        ops,
		},
	)
	assert.Nil(t, err)
	optionsRaw := `{"from":"0x9670d6977d0b10130E5d4916c9134363281B6B0e", "to":"0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674", "data":"0xa9059cbb000000000000000000000000705f9ae78b11a3ed5080c053fa4fa0c52359c674000000000000000000000000000000000000000000000000000000174876e800", "token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A", "value": "0x0"}`
	var options rosettaGethClient.Options
	assert.NoError(t, json.Unmarshal([]byte(optionsRaw), &options))
	assert.Equal(t, &rosettaTypes.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, &options),
	}, preprocessResponse)
}

func TestPreprocessGovernanceDelegate(t *testing.T) {
	networkID := &rosettaTypes.NetworkIdentifier{
		Blockchain: config.DefaultBlockchain,
		Network:    "testnet",
	}
	cfg := &configuration.Configuration{
		Mode:        configuration.ModeOnline,
		Network:     networkID,
		ChainConfig: params.TestChainConfig,
	}
	servicer := construction.NewAPIService(cfg, utils.LoadTypes(), rosettaGethTypes.Errors, nil)
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
]`

	var ops []*rosettaTypes.Operation
	assert.NoError(t, json.Unmarshal([]byte(intent), &ops))
	preprocessResponse, err := servicer.ConstructionPreprocess(
		ctx,
		&rosettaTypes.ConstructionPreprocessRequest{
			NetworkIdentifier: networkID,
			Operations:        ops,
		},
	)
	assert.Nil(t, err)
	optionsRaw := `{"from":"0x9670d6977d0b10130E5d4916c9134363281B6B0e", "to":"0x705f9aE78b11a3ED5080c053Fa4Fa0c52359c674", "data":"0x5c19a95c000000000000000000000000705f9aE78b11a3ED5080c053Fa4Fa0c52359c674", "token_address":"0xF8B089026CaD7DDD8CB8d79036A1ff1d4233d64A", "value": "0x0"}`
	var options rosettaGethClient.Options
	assert.NoError(t, json.Unmarshal([]byte(optionsRaw), &options))
	assert.Equal(t, &rosettaTypes.ConstructionPreprocessResponse{
		Options: forceMarshalMap(t, &options),
	}, preprocessResponse)
}

// marshalJSONMap converts an interface into a map[string]interface{}.
func marshalJSONMap(i interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func forceMarshalMap(t *testing.T, i interface{}) map[string]interface{} {
	m, err := marshalJSONMap(i)
	if err != nil {
		t.Fatalf("could not marshal map %s", rosettaTypes.PrintStruct(i))
	}

	return m
}
