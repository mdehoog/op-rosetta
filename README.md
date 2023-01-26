<div align="center">
  <br />
  <br />
  <a href="https://www.rosetta-api.org"><img alt="Rosetta" src="https://www.rosetta-api.org/img/rosetta_header.png" width=600 height=100></a>
  <br />
  <h3><a href="https://github.com/mdehoog/op-rosetta">op-rosetta</a> is a Rosetta-compatible API for Optimism Bedrock, a low-cost and lightning-fast Ethereum L2 blockchain. It is built on <a href="https://github.com/coinbase/rosetta-geth-sdk">coinbase/rosetta-geth-sdk</a>.</h3>
  <br />
</div>

## Overview

`op-rosetta` provides an executable [client](./app/client.go) that plugs into the [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk). This provides Rosetta API compatibility for Optimism Bedrock, a low-cost and lightning-fast Ethereum L2 blockchain.

To learn more about the Rosetta API, you can find more online at [rosetta-api.org](https://www.rosetta-api.org/).


### Running

To run the `op-rosetta` client, you can use the following command:

```
make build
export CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'
MODE=ONLINE PORT=8080 NETWORK=mainnet GETH=https://mainnet.optimism.io bin/op-rosetta
```

### Testing with `rosetta-cli`

_NOTE: `op-rosetta` must be running on the specified host and port provided in the configuration file. For local testing, this can be done as described in the [Running](#running) section, which will run an instance on localhost, port 8080._

See [configs/README.md](./configs/README.md) for more information.

### License

All files within this repository, including code adapted from [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk), is licensed under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0) unless explicitly stated otherwise.
