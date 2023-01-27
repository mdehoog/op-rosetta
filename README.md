<img align="right" width="150" height="150" top="100" src="./assets/optimism.png">

# op-rosetta • [![tests](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml/badge.svg?label=tests)](https://github.com/mdehoog/op-rosetta/actions/workflows/unit-tests.yaml) [![lints](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml/badge.svg)](https://github.com/mdehoog/op-rosetta/actions/workflows/lints.yaml) 

> **Warning**
> WIP; Not bedrock-compatible yet.

Provides Rosetta API compatibility for Optimism Bedrock, a low-cost and lightning-fast Ethereum L2 blockchain.

## Overview

`op-rosetta` is an executable [client](./app/client.go) extending the [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk).

To learn more about the Rosetta API, you can find more online at [rosetta-api.org](https://www.rosetta-api.org/).

### Running

Run the `op-rosetta` client for goerli like so: `make run-optimism-goerli`.

A mainnet client can be run with `make run-optimism-mainnet`.

> **Note**
> Mainnet will only be supported once Optimism mainnet is upgraded to its Bedrock Release.


### Testing with `rosetta-cli`

_NOTE: `op-rosetta` must be running on the specified host and port provided in the configuration file. For local testing, this can be done as described in the [Running](#running) section, which will run an instance on localhost, port 8080._

See [configs/README.md](./configs/README.md) for more information.

### License

All files within this repository, including code adapted from [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk), is licensed under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0) unless explicitly stated otherwise.
