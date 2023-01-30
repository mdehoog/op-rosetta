package client

import (
	"fmt"

	"github.com/coinbase/rosetta-geth-sdk/configuration"

	evmClient "github.com/coinbase/rosetta-geth-sdk/client"
)

// OpClient wraps the [evmClient.SDKClient] to add Optimism-specific functionality.
type OpClient struct {
	*evmClient.SDKClient
}

// NewOpClient creates a client that can interact with the Optimism network.
func NewOpClient(cfg *configuration.Configuration) (*OpClient, error) {
	client, err := evmClient.NewClient(cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize client: %w", err)
	}
	return &OpClient{client}, nil
}
