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

```
make build
export CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'
MODE=ONLINE PORT=5000 NETWORK=mainnet GETH=https://mainnet.optimism.io bin/op-rosetta
```

### Testing with `rosetta-cli`

To validate `rosetta-ethereum`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/testnet/config.json` - This command validates that the Data API implementation is correct using the ethereum `testnet` node. It also ensures that the implementation does not miss any balance-changing operations.
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/testnet/config.json` - This command validates the Construction API implementation. It also verifies transaction construction, signing, and submissions to the `testnet` network.
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/mainnet/config.json` - This command validates that the Data API implementation is correct using the ethereum `mainnet` node. It also ensures that the implementation does not miss any balance-changing operations.

### License

All files within this repository, including code adapted from [rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk), is licensed under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0) unless explicitly stated otherwise.
