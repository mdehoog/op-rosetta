# op-rosetta

`op-rosetta` is a Rosetta-compatible API for Optimism Bedrock. It is based on [coinbase/rosetta-geth-sdk](https://github.com/coinbase/rosetta-geth-sdk).

### Running

```
make build
export CHAIN_CONFIG='{ "chainId": 10, "terminalTotalDifficultyPassed": true }'
MODE=ONLINE PORT=5000 NETWORK=mainnet GETH=https://mainnet.optimism.io bin/op-rosetta
```
