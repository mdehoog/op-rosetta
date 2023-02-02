<img align="right" width="150" height="150" top="100" src="./assets/rosetta.png">

# op-rosetta • [![tests](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml/badge.svg?label=tests)](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml) [![lints](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml/badge.svg)](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml)

Provides Rosetta API compatibility for Optimism Bedrock, a low-cost and lightning-fast Ethereum L2 blockchain.


## Overview

`op-rosetta` is an executable [client](./app/client.go) extending the [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk) which itself extends [rosetta-sdk-go](https://github.com/coinbase/rosetta-sdk-go).

To learn more about the Rosetta API, you can find more online at [rosetta-api.org](https://www.rosetta-api.org/).


## Setup

First, clone the repository:

```bash
git clone https://github.com/mdehoog/op-rosetta.git
```

Then, configure your own `.env` file by copying the contents of `.env.example` to `.env`.

```bash
cp .env.example .env
```

The `.env` variables will need to be configured.


## Usage

Run the `op-rosetta` client for goerli like so: `make run-optimism-goerli`.

As part of the `run-optimism-goerli` step inside the [makefile](./Makefile), a few parameters are specified before running the `bin/op-rosetta` executable:

```
CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'	\
MODE=ONLINE \
PORT=8080 \
BLOCKCHAIN=Optimism \
NETWORK=Goerli \
ENABLE_TRACE_CACHE=true \
ENABLE_GETH_TRACER=true \
GETH=${OPTIMISM_GOERLI_NODE} \
GENESIS_BLOCK_HASH=${OPTIMISM_GOERLI_GENESIS_BLOCK_HASH} \
bin/op-rosetta
```

These parameters are configured as follows:
- `CHAIN_CONFIG` has the same field values as the genesis config’s. But only the `chainId` and `terminalTotalDifficultyPassed` values are needed.
- Setting `MODE` to `ONLINE` permits outbound connections.
- `PORT` is set to `8080` - the default connection for `rosetta-cli`.
- `BLOCKCHAIN` is `Optimism` and `NETWORK` is `Goerli`
- We enable geth tracing and caching with `ENABLE_TRACE_CACHE` and `ENABLE_GETH_TRACER`
- `GETH` points to a running op-geth *archive* node. This should be set in your `.env` file by setting `OPTIMISM_GOERLI_NODE` to a node url.
- `GENESIS_BLOCK_HASH` is the block hash of the genesis block. This should be set in your `.env` file by setting `OPTIMISM_GOERLI_GENESIS_BLOCK_HASH` to the goerli bedrock genesis block hash.

A mainnet client can be run with `make run-optimism-mainnet`.

> **Note**
> Mainnet will only be supported once Optimism mainnet is upgraded to its Bedrock Release.

The `run-optimism-mainnet` step inside the [makefile](./Makefile) specifies a similar set of variables as above:

```
CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'	\
MODE=ONLINE \
PORT=8080 \
BLOCKCHAIN=Optimism \
NETWORK=Mainnet \
ENABLE_TRACE_CACHE=true \
ENABLE_GETH_TRACER=true \
GETH=${OPTIMISM_MAINNET_NODE} \
GENESIS_BLOCK_HASH=${OPTIMISM_MAINNET_GENESIS_BLOCK_HASH} \
bin/op-rosetta
```

These parameters are configured as follows:
- `CHAIN_CONFIG` has the same field values as the genesis config’s. But only the `chainId` and `terminalTotalDifficultyPassed` values are needed.
- Setting `MODE` to `ONLINE` permits outbound connections.
- `PORT` is set to `8080` - the default connection for `rosetta-cli`.
- `BLOCKCHAIN` is `Optimism` and `NETWORK` is `Mainnet`
- We enable geth tracing and caching with `ENABLE_TRACE_CACHE` and `ENABLE_GETH_TRACER`
- `GETH` points to a running op-geth *archive* node. This should be set in your `.env` file by setting `OPTIMISM_MAINNET_NODE` to a node url.
- `GENESIS_BLOCK_HASH` is the block hash of the genesis block. This should be set in your `.env` file by setting `OPTIMISM_MAINNET_GENESIS_BLOCK_HASH` to the mainnet bedrock genesis block hash.


## Testing

_NOTE: `op-rosetta` must be running on the specified host and port provided in the configuration file. For local testing, this can be done as described in the [Running](#running) section, which will run an instance on localhost, port 8080._

See [configs/README.md](./configs/README.md) for more information.


## Layout

```
├── assets
│   └── rosetta.png
├── cmd
│   └── main.go -> The `op-rosetta` executable entrypoint
├── configs
│   ├── optimism/ -> Optimism Mainnet and Goerli config files
│   └── README.md
└── pkg
    ├── client/ -> Optimism Rosetta API client
    ├── common/ -> Standalone common functionality
    ├── config/ -> Configuration options
    ├── handlers/ -> Rosetta API handlers
    ├── logging/ -> Logging utilities
    └── utils/ -> Package utilities
```

## License

All files within this repository, including code adapted from [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk), is licensed under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0) unless explicitly stated otherwise.
