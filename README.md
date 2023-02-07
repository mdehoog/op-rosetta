<img align="right" width="150" height="150" top="100" src="./assets/rosetta.png">

# op-rosetta • [![tests](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml/badge.svg?label=tests)](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml) [![lints](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml/badge.svg)](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml)

Provides Rosetta API compatibility for Optimism Bedrock, a low-cost and lightning-fast Ethereum L2 blockchain.


## Overview

`op-rosetta` is an executable [client](./app/client.go) extending the [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk) which itself extends [rosetta-sdk-go](https://github.com/coinbase/rosetta-sdk-go).

To learn more about the Rosetta API, you can find more online at [rosetta-api.org](https://www.rosetta-api.org/).


## Run Rosetta Server
1. Run `make build` to build Rosetta Server
2. Configure the following environment variables
```
MODE: ONLINE or OFFLINE
PORT: Rosetta Server port
BLOCKCHAIN: Blockchain (e.g. Optimism)
NETWORK: Network (e.g. Goerli)
FILTER: true (Rosetta Validation validates tokens according to the filter) or false (Rosetta Validation validates all tokens)
GENESIS_BLOCK_HASH: Genesis block hash
GETH: Native node URL
CHAIN_CONFIG: '{ "chainId": <Chain ID>, "terminalTotalDifficultyPassed": true }'
```
3. Run `make run` to run Rosetta Server

## Layout

```
├── assets
│   └── rosetta.png
├── cmd
│   └── main.go -> The `op-rosetta` executable entrypoint
└── pkg
    ├── client/ -> Rosetta API client
    ├── common/ -> Standalone common functionality
    ├── config/ -> Configuration options
    ├── handlers/ -> Rosetta API handlers
    ├── logging/ -> Logging utilities
    └── utils/ -> Package utilities
```

## License

All files within this repository, including code adapted from [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk), is licensed under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0) unless explicitly stated otherwise.
