package client

import (
	"fmt"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
	configuration "github.com/coinbase/rosetta-geth-sdk/configuration"
)

// OpClient wraps the [evmClient.SDKClient] to add Optimism-specific functionality.
//
//go:generate mockery --name OpClient --output ./mocks
type OpClient struct {
	InternalClient
}

// NewOpClient creates a client that can interact with the Optimism network.
func NewOpClient(cfg *configuration.Configuration) (*OpClient, error) {
	client, err := evmClient.NewClient(cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize client: %w", err)
	}
	return &OpClient{client}, nil
}
